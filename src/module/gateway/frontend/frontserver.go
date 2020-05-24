package frontend

import (
	"PushSystem/src/module/gateway"
	"PushSystem/src/pkg/tcpserver"
	"PushSystem/src/pkg/tcpserver/basenet"
	"sync"
)

type FrontModule struct {
	frontEnd tcpserver.TcpServer
	SessionManager * SessionManager
	SessionByIpManager * SessionManagerByIp

}

var _instance * FrontModule
var once sync.Once

func GetInstance() *FrontModule {
	once.Do(func(){
		_instance = &FrontModule {}
	})
	return _instance
}

/*
   初始化
*/
func (handle *FrontModule) Init(){
	handle.frontEnd = &basenet.NetServer{} //采用go程方式的结构可以改为epoll方式
	handle.SessionManager =  GetFrontSessionInstance()
	handle.SessionByIpManager =  GetFrontSessionByIpInstance()

	handle.frontEnd.SetCallBackHandle(handle)
	handle.frontEnd.SetAcceptCallback(gateway.FrontOnAccept)
	handle.frontEnd.SetReceiveCallback(gateway.FrontOnReceive)
	handle.frontEnd.SetCloseCallback(gateway.FrontOnClose)
}

/**/
func (handle *FrontModule) Start(config gateway.GateWayConfig){
	handle.frontEnd.Create(config.Frontend.Ip , config.Frontend.Port)
}

//
func (handle *FrontModule) Stop(){
	handle.frontEnd.ShutDown()
}


