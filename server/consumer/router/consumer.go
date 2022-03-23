package router

import (
	"log"
	"svc-whatsapp/server/consumer/handlers"
)

type Consumer struct {
	handlers.Handler
}

func NewConsumer(handler handlers.Handler) Consumer {
	return Consumer{handler}
}

func (consumer Consumer) Register() {

	cons := consumer.Contract.NsqConsumer
	whatsappHandler := handlers.NewWhatsappConsumerHandler(consumer.Handler)

	cons.Handle("send.message", "channel1", func(message []byte) {
		log.Printf("send.message >> %v\n", string(message))
		// Whatsapp consume route
		whatsappHandler.ConsumeSendMessage(message)
	})

	//cons.Handle("recv.message", "channel1", func(message []byte) {
	//	fmt.Printf("recv.message >> %v\n", string(message))
	//})

	cons.StartListening()
}
