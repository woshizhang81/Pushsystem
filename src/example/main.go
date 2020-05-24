package main

import (
	config2 "PushSystem/src/config"
	"PushSystem/src/module"
	"PushSystem/src/module/gateway"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init(){
	//加载配置
	//创建goroutine pool
	//redis pool

}

var thApp module.Module
func main(){

	appConfig := config2.GetInstance().LoadConfig()
	if nil ==  appConfig {
		os.Exit(1)
	}

	fmt.Println("app name ",appConfig.AppName, "is starting")
	if appConfig.Role== "gateway" {
		thApp = &gateway.GateWay{}
	}else {
		fmt.Println("Role",appConfig.Role, "is not supported")
		os.Exit(1)
	}
	thApp.Start()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	// Stop the service gracefully.
	thApp.Stop()
}



