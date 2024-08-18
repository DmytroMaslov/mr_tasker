package main

import (
	"fmt"
	"log"
	"mr-tasker/cmd/mrtasker/app"
	"mr-tasker/configs"
)

func main() {
	fmt.Println("run app")
	config, err := configs.GetConfig()
	if err != nil {
		log.Fatal("failed to ")
	}
	app.Run(config)
}
