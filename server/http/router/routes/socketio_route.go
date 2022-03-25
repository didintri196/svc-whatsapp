package routes

import (
	"svc-whatsapp/server/http/handlers"
	handler "svc-whatsapp/server/http/handlers/helper"

	"github.com/gin-gonic/gin"

	socketio "github.com/googollee/go-socket.io"
)

type WhatsappSocketioRoute struct {
	Socket  *socketio.Server
	Handler handler.Handler
	Ctx     *gin.Context
}

func NewWhatsappSocketioRoute(so *socketio.Server, handler handler.Handler) WhatsappSocketioRoute {
	return WhatsappSocketioRoute{Socket: so, Handler: handler}
}

func (route WhatsappSocketioRoute) RegisterRoute() {

	// handlers
	whatsappSocketHandler := handlers.NewWhatsappSocketHandler(route.Handler)

	// Users Route
	route.Socket.OnConnect("/", whatsappSocketHandler.OnConnect)
	route.Socket.OnError("/", whatsappSocketHandler.OnError)
	route.Socket.OnDisconnect("/", whatsappSocketHandler.OnDisconnect)

	route.Socket.OnEvent("/", "notice", whatsappSocketHandler.EventNotice)
	route.Socket.OnEvent("/", "generate_qrcode", whatsappSocketHandler.EventReqQrcode)
}
