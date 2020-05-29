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
//	"github.com/letsfire/factory"
//	"time"
	"os"
	"container/list"
	"Pushsystem/src/const"
)
const HbDur = 3 //s

type FrontModule struct {
	ChanArray [_const.SlotNum]chan uint8
	frontEnd tcpserver.TcpServer
	SessionManager * SessionManager
	SessionByIpManager * SessionManagerByIp
	//SlotPool  *factory.Master
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

func process(handle interface{},id int ,c  <-chan uint8){
	obj := handle.(*FrontModule)
	for {
		select {
		case _, ok := <-c :
			if ok {
				//fmt.Println("hbcheck",id,time.Now().Unix())
				obj.SessionManager.HBCheckBySlot(id,HbDur)
			} else {
				//收到关闭信号 退出 go程
				break
			}
		}
	}
}

func (handle *FrontModule) CreateGoPool(){
	for i:=0 ;i< _const.SlotNum ; i++ {
		c := make(chan uint8)
		handle.ChanArray[i] =  c  //保存生成的chan
		go  process(handle ,i  , c )
	}
}

func (handle *FrontModule) HBCheckNotify(){
	for i:=0 ;i< _const.SlotNum ; i++ {
		handle.ChanArray[i] <- 1  //循环通知所有的go程 开始执行心跳检测
	}
}

func (handle *FrontModule) DestroyGoPool(){
	for i:=0 ;i < _const.SlotNum ; i++ {
		close(handle.ChanArray[i])
	}
}

/* 初始化 */
func (handle *FrontModule) Init(){
	handle.frontEnd = &basenet.NetServer{} //采用go程方式的结构可以改为epoll方式
	handle.SessionManager =  GetFrontSessionInstance()
	handle.SessionByIpManager =  GetFrontSessionByIpInstance()
	handle.Channel = &channel.UpStreamChannel{}
	handle.Channel.Init()

	handle.CreateGoPool()

	//handle.SlotPool = factory.NewMaster(_const.SlotNum,_const.SlotNum) //同slot数目相同
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
	handle.DestroyGoPool()
}

func FrontOnAccept (handle interface{} ,conn net.Conn){
	module := handle.(*FrontModule)
	ipAddr := conn.RemoteAddr().String()
	fmt.Println("New Client Connected",ipAddr)

	session := SessionByIp{}
	session.Init()
	session.Conn = conn
	module.SessionByIpManager.Add(ipAddr,session)
}

func FrontOnReceive (handle interface{} ,conn net.Conn ,data []byte){
	module := handle.(*FrontModule)
	ipAddr := conn.RemoteAddr().String()
	//fmt.Println("Client Recieve",ipAddr,"msg:",string(data[:]))
	sessip,ok := module.SessionByIpManager.Get(ipAddr)
	if !ok {
		fmt.Println("fatel error should not be here\n")
		os.Exit(1)
	}

	unit := sessip.(SessionByIp)
	unit.FrameCount ++
	//module.Channel.PutMessage(data[:len(data)])

	listFrame := list.New()
	unit.ProtoCheck.CheckAndGetProtocolBuffer(data,listFrame)
	for item := listFrame.Front();nil != item ;item = item.Next() {
		module.Channel.PutMessage(item.Value.([]byte))
	}
}

/*客户端检测断开*/
func FrontOnClose (handle interface{},conn net.Conn){
	module := handle.(*FrontModule)
	ipAddr := conn.RemoteAddr().String()
	fmt.Println("client close",ipAddr)
	v,ok := module.SessionByIpManager.Get(ipAddr)
	var deviceId 	string
	var deviceType 	DeviceIdType
	if ok {
		deviceId   = v.(SessionByIp).DeviceId
		deviceType = v.(SessionByIp).deviceType
	}else {
		fmt.Println("fatel error should not be here")
		os.Exit(1)
	}
	//同时删除对应客户端Session信息
	module.SessionByIpManager.Delete(ipAddr)
	uniqueId := utils.UniqueId(int32(deviceType),deviceId)
	module.SessionManager.Delete(uniqueId)
}


