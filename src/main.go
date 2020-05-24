package main

import (
	"PushSystem/src/module"
	"PushSystem/src/pkg/tcpserver/basenet"
	"fmt"
	"time"
)
import _ "PushSystem/src/pkg/tcpserver"

var moduleApp module.Module
func init(){
	//moduleApp =
	//加载配置
	//创建goroutine pool
	//redis pool

}

func createServer(server basenet.ServerInterface){
	server.Create("127.0.0.1",8080)
}

func main(){
	server := basenet.NetServer{}
	createServer(&server)
	//server.Init("127.0.0.1",8080)
	fmt.Println("hello")
	for(true){
		time.Sleep(12)
	}
}



