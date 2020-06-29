package frontend

import (
	"Pushsystem/src/pkg/tcpclient"
	"Pushsystem/src/protocol"
)
/*
	设备注册 解析逻辑
*/
func (vServer * VirtualServer) OnGetDeviceRegister(client *tcpclient.TcpClient , protocolData *protocol.Protocol) {

}

func (vServer * VirtualServer) OnGetDeviceHeartBeat(client *tcpclient.TcpClient , protocolData *protocol.Protocol) {

}

/*manager 心跳返回的解析逻辑*/
func (vServer * VirtualServer) OnGetManagerHeartBeat(client *tcpclient.TcpClient, protocolData *protocol.Protocol) {

}

/*manager register返回的解析逻辑*/
func (vServer * VirtualServer) OnGetManagerRegister(client *tcpclient.TcpClient, protocolData *protocol.Protocol) {

}

/*获得push帧的解析逻辑*/
func (vServer * VirtualServer) OnGetPush(client *tcpclient.TcpClient , protocolData *protocol.Protocol) {

}

/*manager register返回的解析逻辑*/
func (vServer * VirtualServer) OnGetTransmission(client *tcpclient.TcpClient , protocolData *protocol.Protocol) {

}
