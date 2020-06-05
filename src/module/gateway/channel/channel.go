package channel

import (
	"github.com/letsfire/factory"
	"fmt"
	"Pushsystem/src/module/gateway/backend"
)

type StreamChannel interface {
	Init()     	//初始化数据流通道
	Start()		//开启通道
	PutMessage(msg interface{}) //添加消息
	Stop()		//停止通道
}

const UpStreamGoRoutineMaxNum 	= 20000  //最大20000个go程池
const UpStreamGoRoutineInitNum 	= 50	   //初始化50个go程大小
const MessageQueueCapity 		= 100000	   //十万个大小队列

/*
	上行数据通道
*/

type UpStreamChannel struct {
	MsgQueue  	chan []byte	//  消息帧
	GoPool  	*factory.Master //
	TaskLine  	*factory.Line
	Consumer	interface{}
}

/*
	go程池 流水线 任务回掉
*/
type MsgUpStream struct {
	Consumer	interface{}
	args interface{}
}

func TaskLineExcute(args interface{}) {
	body := args.(MsgUpStream)
	fmt.Println("将要被发送到后端的数据包",len(body.args.([]byte)),body.args.([]byte))
	backModule := body.Consumer.(*backend.BackModule)
	buf := body.args.([]byte)
	backModule.SendToManager(buf)
	//fmt.Printf("%v", *backEnd)
	// 1. 按轮训方式或者配置的
	// 2. manager权重 发送
}

func (obj * UpStreamChannel) Init()  {
	obj.GoPool = factory.NewMaster(UpStreamGoRoutineMaxNum, UpStreamGoRoutineInitNum)
	obj.MsgQueue = make(chan []byte, MessageQueueCapity)
	obj.TaskLine = obj.GoPool.AddLine(TaskLineExcute)
	obj.Consumer = backend.GetInstance()
}

func (obj * UpStreamChannel) Start()  {

}

func (obj * UpStreamChannel) PutMessage(msg []byte)  {

	par := MsgUpStream{Consumer:obj.Consumer , args:msg}
	obj.TaskLine.Submit(par)
}

func (obj * UpStreamChannel) Stop()  {
	close(obj.MsgQueue)
	obj.GoPool.Shutdown()
}
