package frontend

import "Pushsystem/src/pkg/tcpclient"
import "Pushsystem/src/protocol"

/*manager register返回的解析逻辑*/
func (vServer * VirtualServer) OnGetTransmission(client *tcpclient.TcpClient , protocolData *protocol.Protocol) {
	//该帧为转发帧，是按照 modID 通过kafka 发送到对应的业务模块中
}
