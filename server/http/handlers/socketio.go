package handlers

import (
	"context"
	"fmt"
	"log"
	"svc-whatsapp/domain/requests"
	handler "svc-whatsapp/server/http/handlers/helper"
	"svc-whatsapp/usecase"

	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	waLog "go.mau.fi/whatsmeow/util/log"

	socketio "github.com/googollee/go-socket.io"
)

type WhatsappSocketHandler struct {
	handler.Handler
}

type MyClient struct {
	Contract       *usecase.Contract
	WAClient       *whatsmeow.Client
	SocketConn     socketio.Conn
	Jid            types.JID
	eventHandlerID uint32
	Uuid           string
}

func NewWhatsappSocketHandler(handler handler.Handler) WhatsappSocketHandler {
	return WhatsappSocketHandler{handler}
}

func (handler WhatsappSocketHandler) OnConnect(s socketio.Conn) error {
	log.Println("connected:", s.ID())
	s.SetContext(s.ID())
	s.Join(s.ID())
	s.Emit("reply", "hello world")
	return nil
}

func (handler WhatsappSocketHandler) OnDisconnect(s socketio.Conn, reason string) {
	log.Println("closed", s.Context())
	//log.Println("closed contract", handlers.Contract.ID)
}

func (handler WhatsappSocketHandler) OnError(s socketio.Conn, e error) {
	log.Println("meet error:", e)
}

func (handler WhatsappSocketHandler) EventNotice(s socketio.Conn, msg string) {
	log.Println("notice:", msg)
	s.Emit("reply", "have "+msg)
}

func (mycli *MyClient) register() {
	mycli.eventHandlerID = mycli.WAClient.AddEventHandler(mycli.myEventHandler)
}

func (mycli *MyClient) myEventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.AppStateSyncComplete:
		fmt.Println("Success sync", v.Name)
		deviceID, _ := mycli.Jid.Value()
		fmt.Println("Register", deviceID, "in account uuid ->", mycli.Uuid)
		uc := usecase.NewMDevicesUsecase(mycli.Contract)
		err := uc.AddMDevices(&requests.MDevicesRequest{
			MUserId: mycli.Uuid,
			Jid:     fmt.Sprintf("%s", mycli.Jid),
			Server:  mycli.Jid.Server,
			Phone:   mycli.Jid.User,
		})
		if err != nil {
			fmt.Println("Add Error:", err)
		}
		mycli.SocketConn.Emit("notif", "success_sync")
		mycli.WAClient.Disconnect()
	case *events.PairSuccess:
		mycli.SocketConn.Emit("notif", "success_pair")
		//set JID
		mycli.Jid = v.ID
	}
}

func (handler WhatsappSocketHandler) EventReqQrcode(s socketio.Conn) {
	url := s.URL()
	device := handler.Contract.StoreContainer.NewDevice()
	cli := whatsmeow.NewClient(device, waLog.Stdout("Client", "DEBUG", true))
	store.SetOSInfo("integrasi.in", [3]uint32{0, 1, 0})
	MyclientConnect := &MyClient{
		Contract:   handler.Handler.Contract,
		WAClient:   cli,
		SocketConn: s,
		Uuid:       url.Query().Get("hex"),
	}
	MyclientConnect.register()
	qrChan, err := cli.GetQRChannel(context.Background())
	err = cli.Connect()
	if err != nil {
	} else {
		go func() {
			for evt := range qrChan {
				fmt.Println("Login event:", evt.Event)
				if evt.Event == "code" {
					s.Emit("qrcode", evt.Code)
				} else if evt.Event == "timeout" {
					cli.Disconnect()
					s.Emit("notif", "timeout_scan")
				} else if evt.Event == "success" {
					s.Emit("notif", "success_scan")
				}
			}
		}()
	}
}
