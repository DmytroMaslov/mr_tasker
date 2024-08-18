package app

import (
	"context"
	"fmt"
	"log"
	"mr-tasker/api/native"
	"mr-tasker/api/native/handlers"
	"mr-tasker/configs"
	"mr-tasker/internal/services/user"
	"mr-tasker/internal/storage/dynamodb"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *configs.Config) {
	log.Printf("Hello\n")
	ctx := context.Background()

	dynamodbStorage, err := dynamodb.NewUserStorage(ctx, cfg.Aws)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to create user storage, err: %s", err))
	}
	userService := user.NewUserService(dynamodbStorage)

	cloudHandlers := handlers.NewCloudCrudHandler(userService)

	server := native.NewHttpServer(8080, cloudHandlers.GetHandlers())

	err = server.Serve()
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to serve, err: %s", err))
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	// shout down logic
	err = server.Shutdown(context.Background())
	if err != nil {
		log.Printf("failed to stop server, err: %s", err)
	}

	log.Printf("See you\n")
}
