package backend

import (
	"Pushsystem/src/pkg/tcpserver"
	"sync"
	"Pushsystem/src/pkg/tcpserver/basenet"
	"Pushsystem/src/module/gateway/datadef"
	"Pushsystem/src/module/gateway/callback"
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
	handle.backEnd.SetAcceptCallback (callback.BackOnAccept)
	handle.backEnd.SetReceiveCallback(callback.BackOnReceive)
	handle.backEnd.SetCloseCallback	 (callback.BackOnClose)

}


func (handle *BackModule) Start(config datadef.GateWayConfig){
	handle.backEnd.Create(config.Frontend.Ip,config.Frontend.Port)
}

//
func (handle *BackModule) Stop() {
	handle.backEnd.ShutDown()
}
