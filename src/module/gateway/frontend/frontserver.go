package frontend

import (
	"Pushsystem/src/pkg/tcpserver"
	"Pushsystem/src/pkg/tcpserver/basenet"
	"sync"
	"Pushsystem/src/module/gateway/datadef"
	"Pushsystem/src/module/gateway/channel"
	"net"
	"Pushsystem/src/utils"
	"fmt"
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
	handle.Channel.Init()

	handle.frontEnd.SetCallBackHandle(handle)
	handle.frontEnd.SetAcceptCallback(FrontOnAccept)
	handle.frontEnd.SetReceiveCallback(FrontOnReceive)
	handle.frontEnd.SetCloseCallback(FrontOnClose)
}

func (handle *FrontModule) Start(config datadef.GateWayConfig){
	handle.frontEnd.Create(config.Frontend.Ip , config.Frontend.Port)
}

/*
*/
func (handle *FrontModule) Stop(){
	handle.frontEnd.ShutDown()
}

func FrontOnAccept (handle interface{} ,conn net.Conn){
	module := handle.(*FrontModule)
	ipAddr := conn.RemoteAddr().String()
	fmt.Println("New Client Connected",ipAddr)

	session := SessionByIp{}
	session.Init()
	module.SessionByIpManager.Add(ipAddr,session)
}

func FrontOnReceive (handle interface{} ,conn net.Conn ,data []byte){
	module := handle.(*FrontModule)
	ipAddr := conn.RemoteAddr().String()
	fmt.Println("Client Recieve",ipAddr,"msg:",string(data[:len(data)]))
	session := SessionByIp{}
	session.Init()
	session.Conn = conn
	session.FrameCount ++
	module.Channel.PutMessage(data[:len(data)])
//	var frames [][]byte
//	err := session.ProtoCheck.CheckAndGetProtocolBuffer(data,frames)

/*	if err {
		for _, v := range frames {
			// 将解析出的帧 贴上客户端的ip和端口号
			// 固定第四字节开始50 字节，不够补0
		//	copy(v[3:50],[]byte(ipAddr))
			//调用后端发送到manager里，按机房,qps量加权透传
			//推到后端，由消费者消费 //逻辑就此完成
		}
	}
	module.SessionByIpManager.Add(ipAddr,session)
*/
}

/*客户端检测断开*/
func FrontOnClose (handle interface{},conn net.Conn){
	module := handle.(*FrontModule)
	ipAddr := conn.RemoteAddr().String()
	fmt.Println("client close",ipAddr)
	v,err := module.SessionByIpManager.Get(ipAddr)
	var deviceId string
	var deviceType DeviceIdType
	if err {
		deviceId   = v.(SessionByIp).DeviceId
		deviceType = v.(SessionByIp).DeviceIdType
	}else {
		module.SessionByIpManager.Delete(ipAddr)
		return
	}
	//同时删除对应客户端Session信息
	module.SessionByIpManager.Delete(ipAddr)
	uniqueId := utils.UniqueId(int32(deviceType),deviceId)
	module.SessionManager.Delete(uniqueId)
}


