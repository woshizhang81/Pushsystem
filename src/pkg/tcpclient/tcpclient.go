package tcpclient

import (
	"os"
	"fmt"
	"net"
	"bufio"

	"Pushsystem/src/const"
	"time"
)

type RecvFrameCallBack func  (handle interface{}, buf []byte)
type  TcpClient struct {
	TaskHandle interface{}
	TargetIp   		string
	TargetPort 		uint16
	CreateTime 		uint64
	LastHeartBeatCount 	int64
	HeartBeatErrorTimes uint8  //心跳容错次数 最大为3
	RestartCount 		uint8   //重启次数
	callbackfun 	RecvFrameCallBack
	//IsRecvThreadExist   bool    // 接收线程是否存在 相应逻辑该否补充
	Conn     	*net.TCPConn
	hawkServer  *net.TCPAddr
}


func (client * TcpClient) SetRecieveCallback(handle  interface{},cbfun RecvFrameCallBack ) {
	client.TaskHandle = handle
	client.callbackfun = cbfun
}



/*
	开启定时器
*/


func (client * TcpClient) ReStart() bool {
	//连接服务器
	conn,err := net.DialTCP("tcp",nil,client.hawkServer)
	if err != nil {
		fmt.Printf("connect to hawk server error: [%s]",err.Error())
		return false
	}

	client.Conn = conn
	client.StartRecv()
	return true
}

func (client * TcpClient) Start(server string) bool{
	//开启定时器
	hawkServer,err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		fmt.Printf("hawk server [%s] resolve error: [%s]",server,err.Error())
		os.Exit(1)
	}
	//连接服务器
	conn,err := net.DialTCP("tcp",nil,hawkServer)
	if err != nil {
		fmt.Printf("connect to hawk server error: [%s]",err.Error())
		os.Exit(1)
	}

	client.Conn = conn
	client.hawkServer = hawkServer
	client.StartRecv()
	return true
}

func (client * TcpClient) Send( buf []byte){
	if client.Conn != nil  {
		_ , err := client.Conn.Write(buf)
		if err != nil {
			fmt.Println("send error-->",err.Error())
		}
	}
}

func (client * TcpClient) StartRecv() bool {
	go receivePackets(client)
	return true
}



func (client * TcpClient) Stop() {
	client.Conn.Close() //会退出接收go程
}

func (client * TcpClient) isIpv4() bool {

	return true
}

// 接收数据包
func receivePackets(client *TcpClient) {
	client.RestartCount = 0
	defer client.Conn.Close()
	reader := bufio.NewReader(client.Conn)
	readBufSize := reader.Size()
	buf := make([]byte, readBufSize, readBufSize)   //缓存buffer大小一般为4k
	for {
		reader := bufio.NewReader(client.Conn)
		n , err := reader.Read(buf)
		if n == 0 || err != nil {
			fmt.Printf("read from connect failed, err: %v\n", err)
			break;
		}
		if client.callbackfun != nil  && client.TaskHandle != nil {
			client.callbackfun(client.TaskHandle,buf[:n]) //调用回掉
		}
	}
	fmt.Println("Recieve go routinue finished")
}

