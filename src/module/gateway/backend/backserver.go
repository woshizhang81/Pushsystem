package backend

import (
	"PushSystem/src/module/gateway"
	"PushSystem/src/pkg/tcpserver"
	"PushSystem/src/pkg/tcpserver/basenet"
	"sync"
)

type BackModule struct {
	backEnd tcpserver.TcpServer
}

var _instance * BackModule
var once sync.Once

func GetInstance() *BackModule {
	once.Do(func(){
		_instance = &BackModule {}
	})
	return _instance
}

func (handle *BackModule) Init(){
	handle.backEnd = &basenet.NetServer{} //采用go程方式的结构可以改为epoll方式
	handle.backEnd.SetAcceptCallback(gateway.BackOnAccept)
	handle.backEnd.SetReceiveCallback(gateway.BackOnReceive)
	handle.backEnd.SetCloseCallback(gateway.BackOnClose)
}

/**/
func (handle *BackModule) Start(config gateway.GateWayConfig){
	handle.backEnd.Create(config.Frontend.Ip,config.Frontend.Port)
}

//
func (handle *BackModule) Stop() {
	handle.backEnd.ShutDown()
}
