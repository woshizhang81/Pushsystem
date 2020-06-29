package protocol

import (
	"sync"
	"Pushsystem/src/const"
)

type ProtoFactory struct {

}

func (factory *ProtoFactory) Make(protoType int) *Protocol {
	newProtocol  := &Protocol{}
	switch protoType {
		case _const.ProtoTypeMapRegisterDevice: // 设备注册的协议类型
			 newProtocol.Data = &RegisterBodyDevice{}
			 return newProtocol
		case _const.ProtoTypeMapHeartbeatDevice: // 设备心跳的协议类型
			 return nil
		case _const.ProtoTypeMapRegisterManager: // manager注册的协议类型
			return nil
		case _const.ProtoTypeMapHeartbeatManager: // manager心跳的协议类型
			return nil
		case _const.ProtoTypeMapPush: // 推送的协议类型
			return nil
		case _const.ProtoTypeMapTransmission: // 透传的协议类型
			return nil
		case _const.ProtoTypeMapDisconnect: //客户端主动断开连接协议类型
			return nil
	default:
		return nil
	}
}

var _ProtoFactoryInstence *ProtoFactory
var once sync.Once
func GetInstance() *ProtoFactory {
	once.Do(func(){
		_ProtoFactoryInstence = &ProtoFactory{}
	})
	return _ProtoFactoryInstence
}
