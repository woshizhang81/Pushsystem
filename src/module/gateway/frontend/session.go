package frontend

import (
	"PushSystem/src/utils"
	"net"
	"sync"
)

const slotNum = 500

type DeviceIdType int32

const (
	typeUUID	DeviceIdType = 0
	typeCUID	DeviceIdType = 1
)

type Session struct {
	//ProtoCheck  *protocol.ProtoCheck //负责粘包处理
	DeviceId	string 		// 终端id  推送的唯一标识
	IdType		DeviceIdType //id 类型
	Connection net.Conn		// 连接handle
	State       bool		// 当前状态
	RegisterTime int32      // 注册时间
	HbTimeCount int32		// 心跳时间戳
}

/*
	获取客户端的ip和端口
*/
func (session *Session) remoteAddr() string {
	addr := session.Connection.RemoteAddr()
	return addr.String()
}

func (session *Session) uniqueId() string {
	return utils.UniqueId(int32(session.IdType),session.DeviceId)
}

/*
	前端会话管理类
*/
var sessionManagerInstance *SessionManager
var	sessionOnce sync.Once

func GetFrontSessionInstance() *SessionManager {
	sessionOnce.Do(func(){
		sessionManagerInstance = &SessionManager{}
	})
	return sessionManagerInstance
}

type SessionManager struct{
	syncMapArray [slotNum]sync.Map
	//500个slot 每一个绑定一个sync map 方便心跳 多go程遍历 提高效率
}

func (handle * SessionManager) Add (uniqueId string,session Session) {
	hashcode := utils.HasCode(uniqueId)
	slot := hashcode % slotNum
	handle.syncMapArray[slot].Store(uniqueId,session)
}

func (handle * SessionManager) Get (uniqueId string ,session * Session) (interface{} ,bool) {
	hashcode := utils.HasCode(uniqueId)
	slot := hashcode % slotNum
	return handle.syncMapArray[slot].Load(uniqueId)
}

func (handle * SessionManager) Delete(uniqueId string) {
	hashcode := utils.HasCode(uniqueId)
	slot := hashcode % slotNum
	handle.syncMapArray[slot].Delete(uniqueId)
}


