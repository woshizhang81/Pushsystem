package tcpserver

import "net"

/*
	定义网络模型的公共接口
*/

type OnAccept func(handle interface{},conn net.Conn)
type OnReceive func( handle interface{} ,conn net.Conn ,data []byte)
type OnClose  func(handle interface{},conn net.Conn)

type TcpServer interface {
	Create(ipAddr string, port uint16) bool
	ShutDown()
	SetCallBackHandle(handle interface{})
	SetAcceptCallback(fun OnAccept)
	SetReceiveCallback(fun OnReceive)
	SetCloseCallback(fun OnClose)
}


