package callback

import (
	//"Pushsystem/src/module"
	"Pushsystem/src/module/gateway/frontend"
	"Pushsystem/src/utils"
	"net"
)

//心跳定時器回調函數
func HeartBeatTask (handle interface{} , dur int32,param interface{}){

}

func FrontOnAccept (handle interface{} ,conn net.Conn){
	module := handle.(frontend.FrontModule)
	ipAddr := conn.RemoteAddr().String()
	session := frontend.SessionByIp{}
	session.Init()
	module.SessionByIpManager.Add(ipAddr,session)
}

func FrontOnReceive (handle interface{} ,conn net.Conn ,data []byte){
	module := handle.(frontend.FrontModule)
	ipAddr := conn.RemoteAddr().String()
	session := frontend.SessionByIp{}
	session.Init()
	session.Conn = conn
	session.FrameCount ++
	var frames [][]byte
	err := session.ProtoCheck.CheckAndGetProtocolBuffer(data,frames)
	if err {
		for _, v := range frames {
			// 将解析出的帧 贴上客户端的ip和端口号
			// 固定第四字节开始50 字节，不够补0
			copy(v[3:50],[]byte(ipAddr))
			//调用后端发送到manager里，按机房,qps量加权透传
			//推到后端，由消费者消费 //逻辑就此完成
			module.Channel.PutMessage(v[3:50])
		}
	}
	module.SessionByIpManager.Add(ipAddr,session)
}

/*客户端检测断开*/
func FrontOnClose (handle interface{},conn net.Conn){
	module := handle.(frontend.FrontModule)
	ipAddr := conn.RemoteAddr().String()
	v,err := module.SessionByIpManager.Get(ipAddr)
	var deviceId string
	var deviceType frontend.DeviceIdType
	if err {
		 deviceId   = v.(frontend.SessionByIp).DeviceId
		 deviceType = v.(frontend.SessionByIp).DeviceIdType
	}else {
		module.SessionByIpManager.Delete(ipAddr)
		return
	}
	//同时删除对应客户端Session信息
	module.SessionByIpManager.Delete(ipAddr)
	uniqueId := utils.UniqueId(int32(deviceType),deviceId)
	module.SessionManager.Delete(uniqueId)
}


func BackOnAccept ( handle interface{} ,conn net.Conn){

}

func BackOnReceive (handle interface{} ,conn net.Conn ,data []byte){

}

func BackOnClose  (handel interface{},conn net.Conn){

}


