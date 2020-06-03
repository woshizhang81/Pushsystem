package backend

import (
	"Pushsystem/src/protocol"
	"net"
	"Pushsystem/src/utils"
	"sync"
	"time"
	"math"
	"fmt"
)

type Session struct {
	//ProtoCheck  *protocol.ProtoCheck //负责粘包处理
	ManagerID		[32]byte	// 解析服务器的唯一ID
	ManagerIDC   	uint16		// 解析服务的机房
	Connection 		net.Conn		// 连接handle
	State       	bool		// 当前状态
	RegisterTime 	int64      // 注册时间
	HbTimeCount 	int64		// 心跳时间戳
}
/*
	获取客户端的ip和端口
*/
func (session *Session) remoteAddr() string {
	addr := session.Connection.RemoteAddr()
	return addr.String()
}

func (session *Session) uniqueId() string {
	return utils.UniqueId(int32(session.ManagerIDC),string(session.ManagerID[:]))
}

var sessionManagerInstance *SessionManager
var	sessionOnce sync.Once

func GetFrontSessionInstance() *SessionManager {
	sessionOnce.Do(func(){
		sessionManagerInstance = &SessionManager{}
	})
	return sessionManagerInstance
}

type SessionManager struct{
	//Map sync.Map
	Map SafeMap

}
func (handle * SessionManager) Add (uniqueId string,session *Session) {
	handle.Map.Set(uniqueId,session)
}

func (handle * SessionManager) Get (uniqueId string ) interface{}  {
	return handle.Map.Get(uniqueId)
}

func (handle * SessionManager) Delete(uniqueId string) {
	handle.Map.Delete(uniqueId)
}

func (handle * SessionManager) HBCheckBySlot(slot int, dur int64) {
	handle.Map.Range(func (key,value interface{}) bool{
		uniqueId := key.(string)
		session  := key.(Session)
		curCount := time.Now().Unix()
		if session.HbTimeCount == 0 {
			session.HbTimeCount = curCount
		}else {
			if math.Abs(float64(curCount - session.HbTimeCount)) > float64(dur) {
				fmt.Println("client",uniqueId,"break by hbcheck")
				ipKey := session.remoteAddr()
				obj := GetFrontSessionByIpInstance() //同时
				obj.Delete(ipKey)
			}
		}
		return true
	})
}



type  SessionByIp struct {
	ManagerID		[32]byte	// 解析服务器的唯一ID
	ManagerIDC   	uint16		// 解析服务的机房
	ProtoCheck 		*protocol.ProtoCheck //协议检测
	Conn 			net.Conn
}

func (obj *SessionByIp)Init(){
	obj.ProtoCheck = &protocol.ProtoCheck{}
	obj.ProtoCheck.Init()
}

func (session *SessionByIp) remoteAddr() string {
	addr := session.Conn.RemoteAddr()
	return addr.String()
}

func (session *SessionByIp) uniqueId() string {
	return utils.UniqueId(int32(session.ManagerIDC),string(session.ManagerID[:]))
}

var sessionManagerByIpInstance *SessionManagerByIp
var	sessionByIpOnce sync.Once

func GetFrontSessionByIpInstance() *SessionManagerByIp {
	sessionByIpOnce.Do(func(){
		sessionManagerByIpInstance = &SessionManagerByIp{}
	})
	return sessionManagerByIpInstance
}
type SessionManagerByIp struct{
	Map sync.Map
}

func (handle * SessionManagerByIp) Add (addr string,session SessionByIp) {
	handle.Map.Store(addr,session)
}

func (handle * SessionManagerByIp) Get (addr string ) (interface{} ,bool) {
	return handle.Map.Load(addr)
}

func (handle * SessionManagerByIp) Delete(addr string) {
	handle.Map.Delete(addr)
}
