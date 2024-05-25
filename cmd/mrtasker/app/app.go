package app

import (
	"mr-tasker/api/native"
)

func Run() {
	server := native.NewHttpServer(80)

	server.Serve()
}
