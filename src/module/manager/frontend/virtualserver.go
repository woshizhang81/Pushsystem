package frontend

import (
	"sync"
	"Pushsystem/src/const"
	timer2 "Pushsystem/src/utils/timer"
	"time"
	"fmt"
	"os"
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
	vserver  := handle.(* Link)
	curTimeCount := time.Now().Unix()
	if 	vserver.Client.LastHeartBeatCount != 0 && curTimeCount - vserver.Client.LastHeartBeatCount > _const.ClientHeartBeatCheckDur {
		// 未检测到心跳断开了 此时应该重连
		if vserver.Client.ReStart() {
			vserver.Client.RestartCount ++
			if vserver.Client.RestartCount >= _const.ClientRestartTolerantTimes { //连续重启3次失败
				fmt.Println("连续3次重启失败.... 应该是服务端 关闭了")
				os.Exit(1)
			}
		}
	}
}



type VirtualServer struct {
	LocalAddr	[_const.NetNodeAddrSize]byte	//本地服务的网路地址
	ManagerID   [_const.CommonServerIDSize]byte	//解析服务器唯一ID
	ManagerIDC   uint16		//解析服务的机房
	LinksMap 	sync.Map	//链接映射表  UniqueID 和 对应gateway 链接的关系
	Timer		*timer2.CronTimer  //定时器


}


func (vserver * VirtualServer) Initial() {
	// 读相关 配置文件 填充ManagerID ManagerIDC
	// 初始化kafka客户端
	vserver.Timer = &timer2.CronTimer{}
}

func (vserver * VirtualServer) UnInitial() {

}

func (vserver * VirtualServer) Start() {
	vserver.Timer.Start()
}

func (vserver * VirtualServer) Stop() {
	vserver.Timer.Stop()
}



func (vserver * VirtualServer) CreateTimer(hbdur ,  hbcheck int32) bool {
	//vserver.CreateTimer(_const.ClientHeartBeatDur, _const.ClientHeartBeatCheckDur)
	//开启定时器 30s心跳检测回调
	vserver.Timer = &timer2.CronTimer{}
	vserver.Timer.Start()
	vserver.Timer.Add(HeartBeatSend,vserver, nil , hbdur)
	vserver.Timer.Add(HeartBeatCheck,vserver, nil , hbcheck)
	return true
}

func (vserver * VirtualServer) StopTimer() bool {
	//开启定时器 30s心跳检测回调
	vserver.Timer.Stop()
	return true
}



