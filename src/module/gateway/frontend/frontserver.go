package frontend

import (
	//"Pushsystem/src/module/gateway"
	"Pushsystem/src/pkg/tcpserver"
	"Pushsystem/src/pkg/tcpserver/basenet"
	"sync"
	"Pushsystem/src/module/gateway/datadef"
//	"Pushsystem/src/module/gateway/callback"
	"Pushsystem/src/module/gateway/channel"
//	"Pushsystem/src/module/gateway/callback"
)

type FrontModule struct {
	frontEnd tcpserver.TcpServer
	SessionManager * SessionManager
	SessionByIpManager * SessionManagerByIp
	Channel *channel.UpStreamChannel
}

var _instance * FrontModule
var once sync.Once

func GetInstance() *FrontModule {
	once.Do(func(){
		_instance = &FrontModule {}
	})
	return _instance
}

/* 初始化 */
func (handle *FrontModule) Init(){
	handle.frontEnd = &basenet.NetServer{} //采用go程方式的结构可以改为epoll方式
	handle.SessionManager =  GetFrontSessionInstance()
	handle.SessionByIpManager =  GetFrontSessionByIpInstance()
	handle.Channel = &channel.UpStreamChannel{}

	handle.frontEnd.SetCallBackHandle(handle)
//	handle.frontEnd.SetAcceptCallback(callback.FrontOnAccept)
//	handle.frontEnd.SetReceiveCallback(callback.FrontOnReceive)
//	handle.frontEnd.SetCloseCallback(callback.FrontOnClose)
}

func (handle *FrontModule) Start(config datadef.GateWayConfig){
	handle.frontEnd.Create(config.Frontend.Ip , config.Frontend.Port)
}

/*
*/
func (handle *FrontModule) Stop(){
	handle.frontEnd.ShutDown()
}


