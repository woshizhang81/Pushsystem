package frontend

import (
	"Pushsystem/src/pkg/tcpclient"
	"Pushsystem/src/protocol"
)

/*manager 心跳返回的解析逻辑*/
func (vServer * VirtualServer) OnGetManagerHeartBeat(client *tcpclient.TcpClient , protocolData *protocol.Protocol) {
	// gateway backend 返回的Manager 心跳帧 。
	//todo：暂时不处理,目前利用zookeeper的特性 代替心跳检测方案,有需要再加

}


