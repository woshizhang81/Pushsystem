package tcpclient

import "net"

type OnRecieve func(client * Client, buf []byte)
type Client interface {
	Start(ip string,port uint16) bool
	Send(conn net.Conn , buf []byte)
	StartRecv()
	Close(conn net.Conn)
}
