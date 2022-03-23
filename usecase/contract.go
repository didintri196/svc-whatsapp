package usecase

import (
	"svc-whatsapp/libraries"
	messagebroker "svc-whatsapp/libraries/messagesbroker"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	socketio "github.com/googollee/go-socket.io"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

type Contract struct {
	App            *gin.Engine
	Socket         *socketio.Server
	Validator      *validator.Validate
	StoreContainer *sqlstore.Container
	NsqProducer    messagebroker.Producer
	NsqConsumer    messagebroker.Consumer
	WhatsappWorker *libraries.WorkerPool
}
