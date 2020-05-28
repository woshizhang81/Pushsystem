package backend

import (
	"Pushsystem/src/pkg/tcpserver"
	"sync"
	"Pushsystem/src/pkg/tcpserver/basenet"
	"Pushsystem/src/module/gateway/datadef"
	"net"
	"fmt"
)


type BackModule struct {
	backEnd tcpserver.TcpServer
}

var _instance * BackModule
var once sync.Once

func GetInstance() *BackModule {
	once.Do(func(){
		_instance = &BackModule {}
		_instance.backEnd = &basenet.NetServer{}
		//fmt.Printf("%v",_instance)
	})
	return _instance
}

func (handle *BackModule) Init(){
	//handle.backEnd = &basenet.NetServer{} //采用go程方式的结构可以改为epoll方式
	handle.backEnd.SetAcceptCallback (BackOnAccept)
	handle.backEnd.SetReceiveCallback(BackOnReceive)
	handle.backEnd.SetCloseCallback	 (BackOnClose)

}

func (handle *BackModule) Start(config datadef.GateWayConfig){
	handle.backEnd.Create(config.Frontend.Ip,config.Backend.Port)
}

//
func (handle *BackModule) Stop() {
	handle.backEnd.ShutDown()
}
func BackOnAccept (handle interface{} ,conn net.Conn){
	//module := handle.(*BackModule)
	ipAddr := conn.RemoteAddr().String()
	fmt.Println("New Client Connected",ipAddr)
}

func BackOnReceive (handle interface{} ,conn net.Conn ,data []byte){
	//module := handle.(*BackModule)
	//ipAddr := conn.RemoteAddr().String()

}

/*客户端检测断开*/
func BackOnClose (handle interface{},conn net.Conn){
	// _ := handle.(*BackModule)
	ipAddr := conn.RemoteAddr().String()
	fmt.Println("client close",ipAddr)

}

