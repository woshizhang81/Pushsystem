package channel

import (
	"fmt"
	"github.com/letsfire/factory"
	"Pushsystem/src/module/gateway/frontend"
	"Pushsystem/src/protocol"
	"Pushsystem/src/const"
	"time"
)


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
	frameBuf  :=  body.args.([]byte)
	FrameType :=  protocol.GetTypeFromFrame(frameBuf)
	UniqueID  :=  protocol.GetUniqueIdFromFrame(frameBuf)

	frontEnd := body.Consumer.(*frontend.FrontModule)
	if FrameType == _const.PROTO_TYPE_MAP_REGISTER_DEVICE {
		// 说明是设备注册帧
		// 解包, 更新 SessionByID 和 SessionByIP
		pro := protocol.Protocol{}
		pro.UnPacking(frameBuf)

		//	pro.Data //是否需要
		val,ok := frontEnd.SessionByIpManager.Get(string(pro.Head.ClientAddr[:]))
		if ok {
			//补充 Ip为key的信息
			sess := val.(* frontend.SessionByIp)
			sess.DeviceId = string(pro.Head.DeviceID[:])
			sess.DeviceType = frontend.DeviceIdType(pro.Head.DeviceType)
			frontEnd.SessionByIpManager.Add(string(pro.Head.ClientAddr[:]),sess)

			//添加 UniqueID	w为索引的信息
			newSession := &frontend.Session{}
			newSession.Connection  = sess.Conn
			newSession.DeviceId = string(pro.Head.DeviceID[:])
			newSession.IdType	= frontend.DeviceIdType(pro.Head.DeviceType)
			newSession.RegisterTime = time.Now().Unix()
			newSession.State = false
			frontEnd.SessionByIDManager.Add(UniqueID,newSession)
		}else{
			fmt.Println("逻辑错误异常，此处情况 是，注册帧未返回，连接就断开了。。几乎不可能出现")
			//os.Exit(1)
		}
	} else if FrameType == _const.PROTO_TYPE_MAP_HEARTBEAT_DEVICE {
		//说明是设备心跳帧 更新SessionByID时间戳
		val,ok := frontEnd.SessionByIDManager.Get(UniqueID)
		if ok {
			sess := val.(* frontend.Session)
			sess.HbTimeCount = time.Now().Unix()
			frontEnd.SessionByIDManager.Add(UniqueID,sess)
		}
	} else {
		//只拦截 注册帧和心跳帧，其他直接透传
		//sess := &frontend.Session{}
		val,ok := frontEnd.SessionByIDManager.Get(UniqueID)
		if ok {
			sess := val.(* frontend.Session)
			frontEnd.FrontEnd.Send(sess.Connection , frameBuf )
		}
	}
}

func (obj * DownStreamChannel) Init()  {
	obj.GoPool = factory.NewMaster(_const.DownStreamGoRoutineMaxNum, _const.DownStreamGoRoutineInitNum)
	obj.MsgQueue = make(chan []byte, _const.DownMessageQueueCapity)
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
