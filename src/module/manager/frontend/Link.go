package frontend

import (
	"Pushsystem/src/pkg/tcpclient"
	"Pushsystem/src/protocol"
	"Pushsystem/src/const"
)



type Link struct{
	Client 		 * tcpclient.TcpClient
	GateWayAddr	 [_const.NetNodeAddrSize]byte //要连接的GateWayServer地址
	IsOnline	 bool	//当前链接状态 true 正常,false 离线
	proto   	 *protocol.Protocol //生成注册帧或者其他协议帧的实例
}

