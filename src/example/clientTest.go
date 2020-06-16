package main

import (
	"Pushsystem/src/pkg/tcpclient"
	"os"
	"os/signal"
	"syscall"
	"log"
	"fmt"
)

func main(){
	client := &tcpclient.TcpClient{}
	client.Start("127.0.0.1:8080")

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	client.Stop()
	fmt.Println("nnnnnnnnnnnnnnnnn")
}
