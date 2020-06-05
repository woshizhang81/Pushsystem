package backend

import (
	"Pushsystem/src/pkg/tcpserver"
	"sync"
	"Pushsystem/src/pkg/tcpserver/basenet"
	"Pushsystem/src/module/gateway/datadef"
	"net"
	"fmt"
	"container/list"
	"os"
	"Pushsystem/src/utils"
	"Pushsystem/src/const"
	"Pushsystem/src/module/gateway/channel"
)


type BackModule struct {
	BackEndSer tcpserver.TcpServer
	HeartBeatChan chan uint8
	frontEnd tcpserver.TcpServer
	SessionManager * SessionManager
	SessionByIpManager * SessionManagerByIp
	//SlotPool  *factory.Master
	Channel *channel.DownStreamChannel
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
	handle.BackEndSer = &basenet.NetServer{} //采用go程方式的结构可以改为epoll方式
	handle.BackEndSer.SetAcceptCallback (BackOnAccept)
	handle.BackEndSer.SetReceiveCallback(BackOnReceive)
	handle.BackEndSer.SetCloseCallback	 (BackOnClose)

	handle.SessionManager =  GetFrontSessionInstance()
	handle.SessionByIpManager =  GetFrontSessionByIpInstance()
	handle.Channel = &channel.DownStreamChannel{}
	handle.Channel.Init()

	handle.CreateGoPool()

	handle.BackEndSer.SetCallBackHandle(handle)
}




func process(handle interface{},id int ,c  <-chan uint8){
	obj := handle.(*BackModule)
	for {
		select {
		case _, ok := <-c :
			if ok {
				//fmt.Println("hbcheck",id,time.Now().Unix())
				obj.SessionManager.HBCheckBySlot(id,_const.GateWayFrontHbDur)
			} else {
				//收到关闭信号 退出 go程
				break
			}
		}
	}
}

func (handle *BackModule) CreateGoPool(){
		c := make(chan uint8)
		handle.HeartBeatChan =  c  //保存生成的chan
		go  process(handle ,0  , c )
}

func (handle *BackModule) HBCheckNotify(){
		handle.HeartBeatChan <- 1  //循环通知所有的go程 开始执行心跳检测
}

func (handle *BackModule) DestroyGoPool(){
		close(handle.HeartBeatChan)
}

/* 初始化 */
func (handle *BackModule) Start(config datadef.GateWayConfig){
	handle.frontEnd.Create(config.Frontend.Ip , config.Frontend.Port)
}

func (handle *BackModule) Stop(){
	handle.BackEndSer.ShutDown()
	handle.DestroyGoPool()
}

/*
	向manager发送数据 发送规则 轮训发送
*/
func (handle *BackModule) SendToManager(buf []byte){
	session := handle.SessionManager.Map.GetAverage().(*Session)
	handle.BackEndSer.Send(session.Connection,buf)
}

func BackOnAccept (handle interface{} ,conn net.Conn){
	module := handle.(*BackModule)
	ipAddr := conn.RemoteAddr().String()
	fmt.Println("New Client Connected",ipAddr)

	session := SessionByIp{}
	session.Init()
	session.Conn = conn
	module.SessionByIpManager.Add(ipAddr,session)
}

func BackOnReceive (handle interface{} ,conn net.Conn ,data []byte){
	module := handle.(*BackModule)
	ipAddr := conn.RemoteAddr().String()
	//fmt.Println("Client Recieve",ipAddr,"msg:",string(data[:]))
	sessip,ok := module.SessionByIpManager.Get(ipAddr)
	if !ok {
		fmt.Println("fatel error should not be here\n")
		os.Exit(1)
	}

	unit := sessip.(SessionByIp)

	listFrame := list.New()
	unit.ProtoCheck.CheckAndGetProtocolBuffer(data,listFrame)
	for item := listFrame.Front();nil != item ;item = item.Next() {
		module.Channel.PutMessage(item.Value.([]byte))
	}
}

/*客户端检测断开*/
func BackOnClose (handle interface{},conn net.Conn){
	module := handle.(*BackModule)
	ipAddr := conn.RemoteAddr().String()
	fmt.Println("client close",ipAddr)
	v,ok := module.SessionByIpManager.Get(ipAddr)
	var managerId 	[32]byte
	var managerIDC  uint16
	if ok {
		managerId   = v.(SessionByIp).ManagerID
		managerIDC = v.(SessionByIp).ManagerIDC
	}else {
		fmt.Println("fatel error should not be here")
		os.Exit(1)
	}
	//同时删除对应客户端Session信息
	module.SessionByIpManager.Delete(ipAddr)
	uniqueId := utils.UniqueId(int32(managerIDC),string(managerId[:]))
	module.SessionManager.Delete(uniqueId)
}


