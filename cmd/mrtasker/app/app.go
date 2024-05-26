package app

import (
	"context"
	"fmt"
	"log"
	"mr-tasker/api/native"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	log.Printf("Hello\n")
	server := native.NewHttpServer(80)

	err := server.Serve()
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
