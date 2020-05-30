package main

import (
	timer2 "Pushsystem/src/utils/timer"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func HeartBeatTask (handle interface{} , id int,param interface{}){
	fmt.Println(id,handle,time.Now().Unix())
}

func main() {
	timer := &timer2.CronTimer{}
	timer.Add(HeartBeatTask ,20000, nil , 2)
	timer.Add(HeartBeatTask ,3000000, nil , 10)
	timer.Start()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	// Stop the service gracefully.
	timer.Stop()
	fmt.Println("llllllllllllllllllll")
}
