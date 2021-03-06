package frontend

import (
	"Pushsystem/src/protocol"
	"Pushsystem/src/utils"
	"net"
	"sync"
	"Pushsystem/src/const"
)

type  SessionByIp struct {
	DeviceId 		string 			//注册帧 解析时填充
	DeviceType 	DeviceIdType
	//Idc			 	uint16
	Qps				int		  //qps
	LastFrameCount	uint32    //上次汇聚周期的帧数
	FrameCount		uint32    //IP端口 接收的次数 用于实时计算qps,自增
	ProtoCheck 		*protocol.ProtoCheck //协议检测
	Conn 			net.Conn
}

func (obj *SessionByIp)Init(){
	obj.ProtoCheck = &protocol.ProtoCheck{}
	obj.ProtoCheck.Init()
}

/*
	前端会话管理类
*/
var sessionManagerByIpInstance *SessionManagerByIp
var	sessionByIpOnce sync.Once

func GetFrontSessionByIpInstance() *SessionManagerByIp {
	sessionByIpOnce.Do(func(){
		sessionManagerByIpInstance = &SessionManagerByIp{}
	})
	return sessionManagerByIpInstance
}

type SessionManagerByIp struct{
	syncMapArray [_const.GateWaySlotNum]sync.Map
	//500个slot 每一个绑定一个sync map 方便心跳 多go程遍历 提高效率
}

/*
	addr IP:port
*/
func (handle * SessionManagerByIp) Add (addr string,session * SessionByIp) {
	hashcode := utils.HasCode(addr)
	slot := hashcode % _const.GateWaySlotNum
	handle.syncMapArray[slot].Store(addr,session)
}

func (handle * SessionManagerByIp) Get (addr string ) (interface{} ,bool) {
	hashcode := utils.HasCode(addr)
	slot := hashcode % _const.GateWaySlotNum
	return handle.syncMapArray[slot].Load(addr)
}

func (handle * SessionManagerByIp) Delete(addr string) {
	hashcode := utils.HasCode(addr)
	slot := hashcode % _const.GateWaySlotNum
	handle.syncMapArray[slot].Delete(addr)
}