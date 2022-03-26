package libraries

import (
	"context"
	"fmt"
	"log"
	"svc-whatsapp/domain/models"
	"svc-whatsapp/repositories"
	"svc-whatsapp/utils"

	"go.mau.fi/whatsmeow/types/events"

	"gorm.io/gorm"

	waLog "go.mau.fi/whatsmeow/util/log"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"

	"time"

	"go.mau.fi/whatsmeow/store/sqlstore"
)

type WorkerPool struct {
	Postgres     *gorm.DB
	workersCount int
	jobs         map[string]chan interface{}
	idle         map[string]bool
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

func NewWorkerPool(wcount int, sqlContainer *sqlstore.Container, pg *gorm.DB) *WorkerPool {
	return &WorkerPool{
		Postgres:     pg,
		workersCount: wcount,
		SqlContainer: sqlContainer,
	}
}

func (wp *WorkerPool) generate() (map[string]chan interface{}, map[string]bool) {
	making := make(map[string]chan interface{})
	idle := make(map[string]bool)
	for i := 0; i < wp.workersCount; i++ {
		making[fmt.Sprintf("000-%d", i)] = make(chan interface{}, 0)
		idle[fmt.Sprintf("000-%d", i)] = false
	}

	return making, idle
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
	wp.jobs, wp.idle = wp.generate()

	for i := 0; i < wp.workersCount; i++ {
		// respawn routine
		rouid := fmt.Sprintf("000-%d", i)
		log.Println("Respawn Worker " + rouid)
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

func (wp *WorkerPool) append(items []string, item string) []string {
	var newitems []string
	found := false
	for _, i := range items {
		if i == item {
			found = true
		}
	}

	if !found {
		newitems = append(newitems, item)
	}
	return newitems
}

func (wp *WorkerPool) setIdle(id string, status bool) {
	wp.idle[id] = status
}

func (wp *WorkerPool) worker(ctx context.Context, id string, jobs <-chan interface{}) {
	var client *whatsmeow.Client
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovering from panic in routine %s error is: %v \n", id, r)
		}
	}()

	// set default idle
	wp.setIdle(id, true)

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
					// set worker ID to allocate session
					repo := repositories.NewMDevicesRepository(wp.Postgres)
					dataUpdate := models.Devices{
						WorkerID: id,
					}
					repo.UpdateByJID(wp.Postgres, jobNow.JDID, dataUpdate)

					//connect wamellow
					JDID, _ := types.ParseJID(jobNow.JDID)
					deviceStore, err := wp.SqlContainer.GetDevice(JDID)
					if err != nil {
						panic(err)
					}
					clientLog := waLog.Stdout(id, "DEBUG", true)
					client = whatsmeow.NewClient(deviceStore, clientLog)
					MyclientConnect := &MyClient{
						WAClient:   client,
						WorkerPoll: wp,
						Jid:        JDID,
						UuidWorker: id,
					}
					MyclientConnect.register()
					errClient := client.Connect()
					if errClient != nil {
						fmt.Println("CONNECT CLIENT", errClient)
					}

				case SendMessage:
					if client.IsConnected() {
						fmt.Println(id, "->consume message sendack:", jobNow.Message)
						jID := types.NewJID(jobNow.To, types.DefaultUserServer)
						fmt.Println("RECEIVER MESSAGE JID : ", jID.String())
						msg := &waProto.Message{Conversation: proto.String(jobNow.Message)}
						client.SendMessage(jID, "", msg)
					}
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
		}

	}
}

func (wp *WorkerPool) GetOneIdle() (data string) {
	for key, value := range wp.idle {
		if value {
			return key
		}
	}
	return data
}

func (wp *WorkerPool) GetAllIdle() (data []string) {
	for key, value := range wp.idle {
		if value {
			data = append(data, key)
		}
	}
	return data
}

type MyClient struct {
	WAClient       *whatsmeow.Client
	WorkerPoll     *WorkerPool
	Jid            types.JID
	eventHandlerID uint32
	UuidWorker     string
}

func (mycli *MyClient) register() {
	mycli.eventHandlerID = mycli.WAClient.AddEventHandler(mycli.myEventHandler)
}

func (mycli *MyClient) myEventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.LoggedOut:

		fmt.Println("Success Log Out", mycli.Jid, v.Reason)
		// delete devices because is logged out
		repo := repositories.NewMDevicesRepository(mycli.WorkerPoll.Postgres)
		dataUpdate := models.Devices{
			DeletedAt: time.Now(),
			WorkerID:  "",
		}
		repo.UpdateByJID(mycli.WorkerPoll.Postgres, mycli.Jid.String(), dataUpdate)

		//set worker idle
		mycli.WorkerPoll.setIdle(mycli.UuidWorker, true)
	case *events.ConnectFailureReason:
		fmt.Println("Connect Failure reason :" + v.String())
		if v.IsLoggedOut() {
			fmt.Println("Success Log Out", mycli.Jid)
			// delete devices because is logged out
			repo := repositories.NewMDevicesRepository(mycli.WorkerPoll.Postgres)
			dataUpdate := models.Devices{
				DeletedAt: time.Now(),
				WorkerID:  "",
			}
			repo.UpdateByJID(mycli.WorkerPoll.Postgres, mycli.Jid.String(), dataUpdate)
			//set worker idle
			mycli.WorkerPoll.setIdle(mycli.UuidWorker, true)
		}
	}
}
