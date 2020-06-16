package frontend

import (
	"sync"
	"Pushsystem/src/pkg/tcpclient"
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

type Link struct{
	client 	*tcpclient.TcpClient
	Timer	*timer2.CronTimer
}

func (link * Link) CreateTimer(hbdur ,  hbcheck int32) bool {
	link.CreateTimer(_const.ClientHeartBeatDur, _const.ClientHeartBeatCheckDur)
	//开启定时器 30s心跳检测回调
	link.Timer = &timer2.CronTimer{}
	link.Timer.Start()
	link.Timer.Add(HeartBeatSend,link, nil , hbdur)
	link.Timer.Add(HeartBeatCheck,link, nil , hbcheck)
	return true
}

func (link * Link) StopTimer() bool {
	//开启定时器 30s心跳检测回调
	link.Timer.Stop()
	return true
}
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
	link  := handle.(* Link)
	curTimeCount := time.Now().Unix()
	if 	link.client.LastHeartBeatCount != 0 && curTimeCount - link.client.LastHeartBeatCount > _const.ClientHeartBeatCheckDur {
		// 未检测到心跳断开了 此时应该重连
		if link.client.ReStart() {
			link.client.RestartCount ++
			if link.client.RestartCount >= _const.ClientRestartTolerantTimes { //连续重启3次失败
				fmt.Println("连续3次重启失败.... 应该是服务端 关闭了")
				os.Exit(1)
			}
		}
	}
}



type VirtualServer struct {
	LinksMap sync.Map	//链接映射表  UniqueID 和 对应gateway 链接的关系

}



