package libraries

import (
	"context"
	"fmt"
	"svc-whatsapp/utils"

	waLog "go.mau.fi/whatsmeow/util/log"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"

	"time"

	"go.mau.fi/whatsmeow/store/sqlstore"
)

type WorkerPool struct {
	workersCount int
	jobs         map[string]chan interface{}
	idle         []string
	SqlContainer *sqlstore.Container
}

type SendMessage struct {
	To      string
	Message string
	time    time.Time
}

type ConnectMessage struct {
	JDID string
	time time.Time
}

func NewWorkerPool(wcount int, sqlContainer *sqlstore.Container) *WorkerPool {
	return &WorkerPool{
		workersCount: wcount,
		SqlContainer: sqlContainer,
	}
}

func (wp *WorkerPool) generate() map[string]chan interface{} {
	making := make(map[string]chan interface{})
	for i := 0; i < wp.workersCount; i++ {
		making[fmt.Sprintf("000-%d", i)] = make(chan interface{}, 0)
	}
	return making
}

func (wp *WorkerPool) generateIdle() map[string]bool {
	making := make(map[string]bool)
	for i := 0; i < wp.workersCount; i++ {
		making[fmt.Sprintf("000-%d", i)] = true
	}
	return making
}

func (wp *WorkerPool) Publish(id string, message interface{}) {
	wp.jobs[id] <- message
}

func (wp *WorkerPool) Respawn(ctx context.Context) {
	wp.jobs = wp.generate()

	for i := 0; i < wp.workersCount; i++ {
		// respawn routine
		rouid := fmt.Sprintf("000-%d", i)
		go wp.worker(ctx, rouid, wp.jobs[rouid])
	}

}

func (wp *WorkerPool) remove(items []string, item string) []string {
	var newitems []string

	for _, i := range items {
		if i != item {
			newitems = append(newitems, i)
		}
	}

	return newitems
}

func (wp *WorkerPool) setIdle(id string, status bool) {
	if status {
		wp.idle = append(wp.idle, id)
	} else {
		wp.idle = wp.remove(wp.idle, id)
	}
}

func (wp *WorkerPool) worker(ctx context.Context, id string, jobs <-chan interface{}) {
	var client *whatsmeow.Client
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovering from panic in routine %s error is: %v \n", id, r)
		}
	}()

	for {
		select {
		case job, ok := <-jobs:
			{
				if !ok {
					return
				}
				wp.setIdle(id, false)
				switch jobNow := job.(type) {
				case ConnectMessage:
					fmt.Println(id, "->consume message connectack:", jobNow.JDID)
					JDID, _ := types.ParseJID(jobNow.JDID)
					deviceStore, err := wp.SqlContainer.GetDevice(JDID)
					if err != nil {
						panic(err)
					}
					clientLog := waLog.Stdout(id, "DEBUG", true)
					client = whatsmeow.NewClient(deviceStore, clientLog)
					err = client.Connect()
					if err != nil {
						panic(err)
					}
				case SendMessage:
					fmt.Println(id, "->consume message sendack:", jobNow.Message)
					jID, _ := types.ParseJID("6285895567978@s.whatsapp.net")
					msg := &waProto.Message{Conversation: proto.String(jobNow.Message)}
					client.SendMessage(jID, "", msg)
				default:
					fmt.Println(id, "->consume message notvalid")
				}
			}

		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			if !utils.IsNil(client) {
				client.Disconnect()
			}

			return
		default:
			//todo : cek apakah nil
			if !utils.IsNil(client) {
				if !client.IsConnected() {
					wp.setIdle(id, false)
				}
			} else {
				wp.setIdle(id, false)
			}
		}

	}
}

func (wp *WorkerPool) GetOneIdle() (data string) {
	if len(wp.idle) > 0 {
		return wp.idle[0]
	}
	return data
}
