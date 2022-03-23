package handlers

import (
	"context"
	"fmt"
	"log"
	handler "svc-whatsapp/server/http/handlers/helper"

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
	WAClient       *whatsmeow.Client
	SocketConn     socketio.Conn
	eventHandlerID uint32
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
	log.Println("closed", s.Context().(string))
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
		mycli.SocketConn.Emit("notif", "Success sync")
		mycli.WAClient.Disconnect()
	}
}

func (handler WhatsappSocketHandler) EventReqQrcode(s socketio.Conn) {
	log.Println("MASUK SINI")
	device := handler.Contract.StoreContainer.NewDevice()
	cli := whatsmeow.NewClient(device, waLog.Stdout("Client", "DEBUG", true))
	store.SetOSInfo("integrasi.in", [3]uint32{0, 1, 0})
	MyclientConnect := &MyClient{
		WAClient:   cli,
		SocketConn: s,
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
					s.Emit("notif", "timeout scan")
				} else if evt.Event == "success" {
					s.Emit("notif", "success scan")
				}
			}
		}()
	}
}
