package gateway

import (
	"PushSystem/src/module/gateway/backend"
	"PushSystem/src/module/gateway/frontend"
	"PushSystem/src/utils"
	timer2 "PushSystem/src/utils/timer"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

const ConfigName = "gateway.toml"
type FrontConfig struct {
	QpsLimit int16//qps 限制
	MaxConn int32   //最大客户端连接数字
	SlotNum int16 		//槽位数量
	BlackListNum int16      //黑名单长度
	Ip string				//ip地址
	Port uint16				//端口号
}

type BackConfig struct {
	QpsLimit int16//qps 限制
	MaxConn int32   //最大客户端连接数字
	SlotNum int16 		//槽位数量
	BlackListNum int16      //黑名单长度
	Ip 		string				//ip地址
	Port 	int16				//端口号
}
/*gateway module 配置*/
type GateWayConfig struct{
	Module string
	Frontend FrontConfig
	Backend	 BackConfig
}

type GateWay struct{
	config  GateWayConfig
	frontEnd frontend.FrontModule
	backEnd	 backend.BackModule
	timer    timer2.CronTimer
}

func (handle *GateWay)loadConfig() bool {
	filename := utils.GetConfigPath()+ ConfigName
	_, err := toml.DecodeFile(filename,&handle.config)
	if err != nil{
		msg := "config file :"+ filename +"load failed"
		log.Fatal(msg)
		os.Exit(1)
		return false
	}
	return true
}

func (handle *GateWay) Start()  {
	//读gateway配置
	ret := handle.loadConfig()
	if ret == false {
		os.Exit(1)	//退出程序
	}
	//开启前端
	handle.frontEnd.Init()
	handle.frontEnd.Start(handle.config)

	//开启后端
	handle.backEnd.Init()
	handle.backEnd.Start(handle.config)
	//开启定时器 30s心跳检测回调
	handle.timer.Add(HeartBeatTask ,handle, nil , 30)
	/*  //待功能补全
	//开启定时器 10s gateway 流量检测回调
	handle.timer.Add(HeartBeatTask ,handle, nil , 10)
	//开启定时器 5s 机器检查回调
	handle.timer.Add(HeartBeatTask ,handle, nil , 10)
	*/
	handle.timer.Start()
	//等待添加任務

}

func (handle *GateWay) Stop()  {
	handle.frontEnd.Stop()
	handle.backEnd.Stop()
	handle.timer.Stop()

}
