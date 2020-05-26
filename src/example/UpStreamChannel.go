package main

import (
	"time"
	"Pushsystem/src/module/gateway/channel"
)

func main(){
	obj := channel.UpStreamChannel{}
	obj.Init()
	for {
		msg := "123456"
		obj.PutMessage([]byte(msg))
		time.Sleep(time.Second)
	}
}
