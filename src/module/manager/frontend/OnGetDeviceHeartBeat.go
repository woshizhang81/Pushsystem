package frontend

import (
	"Pushsystem/src/pkg/tcpclient"
	"Pushsystem/src/protocol"
)

func (vServer * VirtualServer) OnGetDeviceHeartBeat(client *tcpclient.TcpClient , protocolData *protocol.Protocol) {
	//1. 更新 redis 中 设备ID过期时间
	//2. 构造心跳返回帧 返回到设备中，需要协议支持
}


