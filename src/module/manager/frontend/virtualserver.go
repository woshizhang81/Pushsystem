package frontend

import (
	"Pushsystem/src/config"
	"Pushsystem/src/const"
	"Pushsystem/src/pkg/tcpclient"
	"Pushsystem/src/pkg/tools/zkclient"
	"Pushsystem/src/protocol"
	"Pushsystem/src/utils"
	timer2 "Pushsystem/src/utils/timer"
	"crypto"
	"crypto/md5"
	"fmt"
	"os"
	"sync"
)

/*
	虚拟的服务器。
	模拟 通过客户端 与各个geteway 后端建立的连接集合
	需求:
	1. 通过统一的接口 动态，实时的获取 geteway 服务端 并实时建立连接 (zookeeper)
	2. 针对每个geteway 服务 应该有重连机制。//应该维护(10)个链接 ，每个链接应该具备心跳检测和重连机制
	3. 应对每个gateway 所有链接进行 上下行数据进行流量统计(pv,qps)等性能指标进行监控,最好有流量控制机制
	4. 需要实现协议解析
		1> 注册协议 绑定（deviceid 和 devicetype信息） 和 gateway 服务器的关系到缓存中
		2> 心跳。   更新（deviceid 和 devicetype信息） 的缓存信息的 过期时间
		3> 业务透传协议。 根据 协议中的mode_id  透传信息到各个业务中 动态定制kafkatopick
		4> 需要满足按 deviceid 和 devicetype 选择发送到gateway的 链接的功能
*/
/*
	ZooKeeper路径结构 定义
    功能:服务器动态上线下感知
	1. gateway 服务 结构定义
	PushSystem/gatway_{idc(uint16)}   value:json {idc:'tx'}       节点类型，永久节点
	child {parentPath}/{ip+port}	  value:json {}  负载参数用于负载均衡
	2. manager 服务 结构定义
	PushSystem/gatway_{idc(uint16)}   value:json {idc:'tx'}       节点类型，永久节点
	child {parentPath}/{ip+port}	  value:json {}  内容待定{负载参数用于负载均衡}
*/
// 发送心跳帧




type VirtualServer struct {
	ZkServers 	[]string
	LocalAddr	[_const.NetNodeAddrSize]byte	//本地服务的网路地址
	ManagerID   [_const.CommonServerIDSize]byte	//解析服务器唯一ID
	ManagerIDC  uint16		//解析服务的机房
	GateWayIDLinksMap 		sync.Map	//目标target IPaddr(gatewayID)和 对应gateway client 客户端链接的关系
	UniqueIDIpMap 			sync.Map	//终端  UniqueID 和 对应gatewayID 链接的关系 存在redis中
	Timer		*timer2.CronTimer  //定时器
	ZkHandel 	*zkclient.ZkClient
}


func (vServer * VirtualServer) Init() {
	zkConfig := config.GetZkInstance().LoadConfig()
	if zkConfig == nil {
		fmt.Println("zookeeper config loading failed")
		os.Exit(1)
	}
	vServer.ZkServers = zkConfig.Servers
	vServer.ZkHandel = &zkclient.ZkClient{ZkAddr:zkConfig.Servers}
	// 读相关 配置文件 填充ManagerID ManagerIDC
	// 初始化kafka客户端
	vServer.Timer = &timer2.CronTimer{}
	vServer.CreateTimer(_const.ClientHeartBeatDur, _const.ClientHeartBeatCheckDur)
}


func (vServer * VirtualServer) Start() {
	if !vServer.ZkHandel.Start() {
		os.Exit(1)
	}
	//添加监控路径  监控gateway 的节点事件
	vServer.ZkHandel.AddPathEvents(_const.ZkGateWayParentNodeName,
		vServer,
		CallBackPathCreated,
		CallBackPathDeleted,
		CallBackPathContextChanged,
		CallBackPathChildNumChanged	)
	vServer.Timer.Start()
}

func (vServer * VirtualServer) Stop() {
	vServer.ZkHandel.Stop()
	vServer.StopTimer()
}

func (vServer * VirtualServer) CreateTimer(hbDur ,  hbCheck int32) bool {
	//开启定时器 30s心跳检测回调
	vServer.Timer.Add(HeartBeatSend  , vServer, nil , hbDur)
	vServer.Timer.Add(HeartBeatCheck , vServer, nil , hbCheck)
	return true
}

func (vServer * VirtualServer) StopTimer() bool {
	//开启定时器 30s心跳检测回调
	vServer.Timer.Stop()
	return true
}

func CallBackPathCreated  (handle interface{},path string ,nodeValue []byte) {

}

func CallBackPathDeleted (handle interface{},path string){

}

func CallBackPathContextChanged (handle interface{}  , path string, latestPathValue []byte ,currentPathValue []byte){

}

func HeartBeatSend (handle interface{}, id int , param interface{}) {
	vServer  := handle.(* VirtualServer)
	//curTimeCount := time.Now().Unix()
	vServer.GateWayIDLinksMap.Range(func (key interface{},value interface{}) bool {
		//remoteServer := key.(string)
		cli := value.(*tcpclient.TcpClient)
		//msg := cli.ProtocolHandler.Package()
		cli.ProtocolHandler.PackType 	 	= _const.ProtoTypeMapHeartbeatManager
		 ok , addr , idc :=utils.GetServerInstance()
		if !ok {
			fmt.Println("should not be here")
			panic("should not be here")
		}
		copy(cli.ProtocolHandler.Head.ManagerID[:]  ,[]byte(utils.MD5(addr))[:])
		cli.ProtocolHandler.Head.ManagerIDC = idc
		//todo :心跳帧待定义
		HeartBeatFrame := cli.ProtocolHandler.Package()
		cli.Send(HeartBeatFrame)

		return true
	})
}

// 心跳检查定时器
func HeartBeatCheck(handle interface{}, id int , param interface{}) {
	//fmt.Println("HeartBeatCheck",time.Now().Second())
	//todo:心跳检查 个人觉得没有必要，利用zookeeper高可用集群作判断 足够了,预留 心跳检测的逻辑。作备用方案
	//vServer  := handle.(* VirtualServer)
	//curTimeCount := time.Now().Unix()
	/*vServer.GateWayIDLinksMap.Range(func (key interface{},value interface{}) bool {
		remoteServer := key.(string)
		cli := value.(*tcpclient.TcpClient)
		if 	cli.LastHeartBeatCount != -2 && curTimeCount - cli.LastHeartBeatCount > _const.ClientHeartBeatCheckDur {
			if cli.ReStart() {
				cli.RestartCount++
				if cli.RestartCount >= _const.ClientRestartTolerantTimes { //连续重启1次失败
					fmt.Println(remoteServer,"连续1次重启失败.... 应该是服务端 关闭了")

					//os.Exit(-1)
				}
			}
		}
		return true
	})*/
}

func OnRecvFrame(handle interface{}, remoteAddr string ,buf []byte) {
	//主要的协议解析逻辑在这里
	vServer := handle.(* VirtualServer)
	protoHandler := &protocol.Protocol{}
	protoHandler.UnPacking(buf) //解析出协议的结构体

	cli,err := vServer.GateWayIDLinksMap.Load(utils.MD5(remoteAddr)) //一定存在该client的 映射关系
	var client *tcpclient.TcpClient
	if err {
		client = cli.(* tcpclient.TcpClient)
	}else{
		fmt.Println("this cannot happened !!")
		panic("this cannot happened")
	}
	switch protoHandler.PackType {
	case _const.ProtoTypeMapRegisterDevice:    //设备注册的协议类型
		vServer.OnGetDeviceRegister(client ,protoHandler)
	case _const.ProtoTypeMapHeartbeatDevice:   //设备心跳的协议类型
		vServer.OnGetDeviceHeartBeat(client , protoHandler)
	case _const.ProtoTypeMapRegisterManager:   //manager注册的协议类型
		vServer.OnGetManagerRegister(client , protoHandler)
	case _const.ProtoTypeMapHeartbeatManager:  //manager心跳的协议类型
		vServer.OnGetManagerHeartBeat(client , protoHandler)
	case _const.ProtoTypeMapPush:              //推送的协议类型
		vServer.OnGetPush(client , protoHandler)
	case _const.ProtoTypeMapTransmission:      //透传的协议类型
		vServer.OnGetTransmission(client , protoHandler)
	case _const.ProtoTypeMapDisconnect:        //客户端主动断开连接协议类型
	default:
		fmt.Println("unrecognized commander!!")
	}
}

func CallBackPathChildNumChanged (handle interface{} , path string, changeType uint8 , ChangedNode string) {
	// 主要 //监控子节点的 事件。动态监控gateway 服务器。适用的事件
	vServer := handle.(* VirtualServer)
	zkPath := path
	targetAddr := ChangedNode
	fmt.Println("addr:",targetAddr,"parentpath:",zkPath,"event:",changeType)
	switch changeType {
	case zkclient.ZKChildAdd: //模拟accept 到客户端
		//建立连接
		cli := tcpclient.TcpClient{TaskHandle:vServer,
			CallbackFun:OnRecvFrame,
			ProtocolHandler:&protocol.Protocol{}}
		ok := cli.Start(targetAddr)
		if ok {
			vServer.GateWayIDLinksMap.Store(utils.MD5(targetAddr),cli)
		}
	case zkclient.ZKChildDel: //模拟检测到客户端断开
		cli,ok := vServer.GateWayIDLinksMap.Load(utils.MD5(targetAddr))
		if ok {
			client := cli.(* tcpclient.TcpClient)
			//vServer.UniqueIDIpMap.Delete()
			client.Stop()
			vServer.GateWayIDLinksMap.Delete(utils.MD5(targetAddr))
			//todo： 清除redis绑定信息，或者等待redis 健自动过期
		}
	default:
		break
	}
}




