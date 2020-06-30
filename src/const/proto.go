package _const

const (
	ProtoTypeMapRegisterDevice   =  0 //设备注册的协议类型
	ProtoTypeMapHeartbeatDevice  =  1 //设备心跳的协议类型
	ProtoTypeMapRegisterManager  =  2 //manager注册的协议类型
	ProtoTypeMapHeartbeatManager =  3 //manager心跳的协议类型
	ProtoTypeMapPush             =  4 //推送的协议类型
	ProtoTypeMapTransmission     =  5 //透传的协议类型
	ProtoTypeMapDisconnect       =  6 //客户端主动断开连接协议类型


	ProtoDefaultFrameSize	= 2048

	ProtoFlagUpStream 		= 0x0001 //上下行数据 1 为上行 0 为下行
	ProtoFlagNeedAck		= 0x0002 //是否需要应答 1.需要应答 0.不需要
	ProtoFlagNeed			= 0x0004 //是否需要geteway拆包 //用于注册帧和心跳帧
	protoFlag
)
