package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"svc-whatsapp/domain"
	handler "svc-whatsapp/server/http/handlers/helper"
	"svc-whatsapp/server/http/router"
	"svc-whatsapp/usecase"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Load Configuration
	config, err := domain.LoadConfiguration()
	if err != nil {
		log.Fatal("Error while load configuration, ", err.Error())
	}

	// Insert Handler Contract
	var handler = handler.NewHandler(&usecase.Contract{
		App:            config.App,
		Socket:         config.Socket,
		Validator:      config.Validator,
		StoreContainer: config.StoreContainer,
		NsqProducer:    config.NsqProducer,
		Postgres:       config.Postgres,
		SecretKey:      []byte(config.SecretKey),
	})

	// Register routes
	router.NewRouter(handler).RegisterRoutes()

	srv := domain.HttpListen(config.App)

	// Listening Http
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("%s\n", err)
		}
	}()

	// Listening Socket
	go func() {
		if err := config.Socket.Serve(); err != nil {
			log.Printf("Socket.io %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	err = config.Socket.Close()
	if err != nil {
		log.Fatal("Socket.io forced to shutdown:", err)
	}

	log.Println("Server exiting")

}
