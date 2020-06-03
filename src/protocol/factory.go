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
		case _const.PROTO_TYPE_MAP_REGISTER_DEVICE:   // 设备注册的协议类型
			 newProtocol.Data = &RegisterBodyDevice{}
			 return newProtocol
		case _const.PROTO_TYPE_MAP_HEARTBEAT_DEVICE: // 设备心跳的协议类型
			 return nil
		case _const.PROTO_TYPE_MAP_REGISTER_MANAGER: // manager注册的协议类型
			return nil
		case _const.PROTO_TYPE_MAP_HEARTBEAT_MANAGER: // manager心跳的协议类型
			return nil
		case _const.PROTO_TYPE_MAP_PUSH:  // 推送的协议类型
			return nil
		case _const.PROTO_TYPE_MAP_TRANSMISSION: // 透传的协议类型
			return nil
		case _const.PROTO_TYPE_MAP_DISCONNECT: //客户端主动断开连接协议类型
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
