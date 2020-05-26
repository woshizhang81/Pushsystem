package entrance

import (
	"Pushsystem/src/module/gateway/backend"
	"Pushsystem/src/module/gateway/frontend"
	"Pushsystem/src/utils"
	timer2 "Pushsystem/src/utils/timer"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"Pushsystem/src/module/gateway/datadef"
)

const ConfigName = "gateway.toml"

type GateWay struct{
	config   datadef.GateWayConfig
	frontEnd frontend.FrontModule
	backEnd  backend.BackModule
	timer    timer2.CronTimer
}

func (handle *GateWay)loadConfig() bool {
	filename := utils.GetConfigPath()+ ConfigName
	_, err := toml.DecodeFile(filename,&handle.config)
	if err != nil{
		msg := "datadef file :"+ filename + "load failed"
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
//	handle.timer.Add(HeartBeatTask ,handle, nil , 30)
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
