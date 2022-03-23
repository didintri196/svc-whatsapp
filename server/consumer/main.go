package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"svc-whatsapp/domain"
	"svc-whatsapp/libraries"
	"svc-whatsapp/server/consumer/handlers"
	"svc-whatsapp/server/consumer/router"
	"svc-whatsapp/usecase"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Load Configuration
	config, err := domain.LoadConfiguration()
	if err != nil {
		log.Fatal("Error while load configuration, ", err.Error())
	}

	// Insert Handler Contract
	var handler = handlers.NewHandler(&usecase.Contract{
		Validator:      config.Validator,
		StoreContainer: config.StoreContainer,
		NsqProducer:    config.NsqProducer,
		NsqConsumer:    config.NsqConsumer,
		WhatsappWorker: config.WhatsappWorker,
	})

	// Start Worker
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	handler.Contract.WhatsappWorker.Respawn(ctx)
	handler.Contract.WhatsappWorker.Publish("000-0", libraries.ConnectMessage{JDID: "6285155075517.0:11@s.whatsapp.net"})

	// Register routes
	router.NewConsumer(handler).Register()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	cancel()
	// The context is used to inform the server it has 5 seconds to finish

	// the request it is currently handling
	log.Println("Server exiting")
}
