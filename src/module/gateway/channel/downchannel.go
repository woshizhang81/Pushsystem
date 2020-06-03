package channel

import (
	"fmt"
	"github.com/letsfire/factory"
	"Pushsystem/src/module/gateway/frontend"
)

const DownStreamGoRoutineMaxNum 	= 2000  //最大20000个go程池
const DownStreamGoRoutineInitNum 	= 50	   //初始化50个go程大小
const DownMessageQueueCapity 		= 100000	   //十万个大小队列

/*
	上行数据通道
*/

type DownStreamChannel struct {
	MsgQueue  	chan []byte	//  消息帧
	GoPool  	*factory.Master //
	TaskLine  	*factory.Line
	Consumer	interface{}
}

/*
	go程池 流水线 任务回掉
*/
type MsgDownStream struct {
	Consumer	interface{}
	args interface{}
}

func TaskLineExcuteDown(args interface{}) {
	body := args.(MsgDownStream)
	fmt.Println("将要被发送到前端端的数据包",len(body.args.([]byte)),body.args.([]byte))
	//按UniqueID 找到conn并发送

	//backEnd := body.Consumer.(*backend.BackModule)
	//fmt.Printf("%v", *backEnd)
	// 1. 按轮训方式或者配置的
	// 2. manager权重 发送
}

func (obj * DownStreamChannel) Init()  {
	obj.GoPool = factory.NewMaster(DownStreamGoRoutineMaxNum, DownStreamGoRoutineInitNum)
	obj.MsgQueue = make(chan []byte, DownMessageQueueCapity)
	obj.TaskLine = obj.GoPool.AddLine(TaskLineExcuteDown)
	obj.Consumer = frontend.GetInstance()
}

func (obj * DownStreamChannel) Start()  {

}

func (obj * DownStreamChannel) PutMessage(msg []byte)  {

	par := MsgDownStream{Consumer:obj.Consumer , args:msg}
	obj.TaskLine.Submit(par)
}

func (obj * DownStreamChannel) Stop()  {
	close(obj.MsgQueue)
	obj.GoPool.Shutdown()
}
