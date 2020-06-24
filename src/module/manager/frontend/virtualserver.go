package frontend

import (
	"Pushsystem/src/pkg/tools/zkclient"
	"sync"
	"Pushsystem/src/const"
	timer2 "Pushsystem/src/utils/timer"
	"time"
	"fmt"
	"os"
	"Pushsystem/src/config"
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
func HeartBeatSend (handle interface{}, id int , param interface{}) {
	//fmt.Println("HeartBeatSend",time.Now().Second())
	//client  := handle.(* TcpClient)
	/*if client.LastHeartBeatCount == 0 {
		client.LastHeartBeatCount = time.Now().Unix()
	}*/
}

// 心跳检查定时器
func HeartBeatCheck(handle interface{}, id int , param interface{}) {
	//fmt.Println("HeartBeatCheck",time.Now().Second())
	vServer  := handle.(* Link)
	curTimeCount := time.Now().Unix()
	if 	vServer.Client.LastHeartBeatCount != 0 && curTimeCount - vServer.Client.LastHeartBeatCount > _const.ClientHeartBeatCheckDur {
		// 未检测到心跳断开了 此时应该重连
		if vServer.Client.ReStart() {
			vServer.Client.RestartCount ++
			if vServer.Client.RestartCount >= _const.ClientRestartTolerantTimes { //连续重启3次失败
				fmt.Println("连续3次重启失败.... 应该是服务端 关闭了")
				os.Exit(1)
			}
		}
	}
}



type VirtualServer struct {
	ZkServers 	[]string
	LocalAddr	[_const.NetNodeAddrSize]byte	//本地服务的网路地址
	ManagerID   [_const.CommonServerIDSize]byte	//解析服务器唯一ID
	ManagerIDC  uint16		//解析服务的机房
	GateWayIDLinksMap 		sync.Map	//目标target IPaddr(gatewayID)和 对应gateway 链接的关系
	UniqueIDIpMap 	sync.Map	//终端  UniqueID 和 对应gatewayID 链接的关系 存在redis中
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



func (vServer * VirtualServer) CreateTimer(hbdur ,  hbcheck int32) bool {
	//开启定时器 30s心跳检测回调
	vServer.Timer.Add(HeartBeatSend  , vServer, nil , hbdur)
	vServer.Timer.Add(HeartBeatCheck , vServer, nil , hbcheck)
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

func CallBackPathChildNumChanged (handle interface{} , path string, changeType uint8 , ChangedNode string) {
	// 主要 适用的事件
	if changeType ==  zkclient.ZKChildAdd {
		//
	}
}




