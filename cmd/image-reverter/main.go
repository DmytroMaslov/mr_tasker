package main

import (
	"fmt"
	"log"
	"mr-tasker/cmd/image-reverter/app"
	"time"
)

func main() {
	// https://dou.ua/
	now := time.Now()
	fmt.Println("start sync convert")
	syncConverter := app.SyncImageConvertor{}
	res, err := syncConverter.Convert("https://dou.ua/")
	if err != err {
		log.Fatal(err)
	}
	fmt.Println(res)
	fmt.Printf("finish sync convert, time takes %d millisecond\n", time.Since(now).Milliseconds())

	now = time.Now()
	fmt.Println("start async convert")
	async := app.AsyncImageConvertor{}
	res, err = async.Convert("https://dou.ua/")
	if err != err {
		log.Fatal(err)
	}
	fmt.Println(res)
	fmt.Printf("finish async convert, time takes %d millisecond\n", time.Since(now).Milliseconds())
}
