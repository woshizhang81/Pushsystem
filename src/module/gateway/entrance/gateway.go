package entrance

import (
	_const "Pushsystem/src/const"
	"Pushsystem/src/module/gateway/backend"
	"Pushsystem/src/module/gateway/frontend"
	"Pushsystem/src/utils"
	timer2 "Pushsystem/src/utils/timer"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"Pushsystem/src/module/gateway/datadef"
	"fmt"
	//"time"
)

const ConfigName = "gateway.toml"

type GateWay struct{
	config   datadef.GateWayConfig
	frontEnd *frontend.FrontModule
	backEnd  *backend.BackModule
	timer    *timer2.CronTimer
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
	handle.frontEnd = frontend.GetInstance()
	handle.frontEnd.Init()
	fmt.Printf("%v" , handle.config)
	handle.frontEnd.Start(handle.config)

	//开启后端
	handle.backEnd = backend.GetInstance()
	handle.backEnd.Init()
	handle.backEnd.Start(handle.config)
	//开启定时器 30s心跳检测回调
	handle.timer = &timer2.CronTimer{}
	handle.timer.Add(GateWayFrontEndHeartBeatTask ,handle, nil , _const.GateWayFrontHbDur)
	handle.timer.Add(HostStateCheck ,handle, nil , _const.GateWayHostStateCheckDur)
	  //待功能补全
	//开启定时器 10s gateway 后端心跳检测回调
	handle.timer.Add(GateWayBackEndHeartBeatTask,handle, nil , _const.GateWayBackHbDur)
	handle.timer.Add(GateWayBackEndLoadBalanceCheck,handle, nil , _const.GateWayBackLoadBalanceDur)
	//前端 流量检测定时回调
	handle.timer.Add(GateWayFrontEndFlowRateCheck,handle, nil , _const.GateWayFrontFlowRateDur)
	handle.timer.Start()
}

func (handle *GateWay) Stop()  {
	handle.frontEnd.Stop()
	handle.backEnd.Stop()
	handle.timer.Stop()
}

//前端 心跳检测
func GateWayFrontEndHeartBeatTask (handle interface{} , id int,param interface{}){
	//fmt.Println(id,handle,time.Now().Unix())
	gateway := handle.(*GateWay)
	gateway.frontEnd.HBCheckNotify()
}

// 10s 检查一次主机状态
func HostStateCheck (handle interface{} , id int,param interface{}) {

}

// 后端心跳检测逻辑
func GateWayBackEndHeartBeatTask(handle interface{} , id int,param interface{}){

}

//前端流量控制逻辑
func GateWayFrontEndFlowRateCheck(handle interface{} , id int,param interface{}){

}

//后端负载均衡定时回调
func GateWayBackEndLoadBalanceCheck(handle interface{} , id int,param interface{}){

}




