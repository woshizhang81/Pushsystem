package tcpclient

import (
	"os"
	"fmt"
	"net"
	"bufio"
	timer2 "Pushsystem/src/utils/timer"
	"Pushsystem/src/const"
	"time"
)

type  BaseClient struct {
	TargetIp   		string
	TargetPort 		uint16
	CreateTime 		uint64
	LastHeartBeatCount 	int64
	HeartBeatErrorTimes uint8  //心跳容错次数 最大为3
	RestartCount 		uint8   //重启次数
	Timer		*timer2.CronTimer
	Conn     	*net.TCPConn
	hawkServer  *net.TCPAddr
	stopChan    chan struct{}
}


// 发送心跳帧
func HeartBeatSend (handle interface{}, id int , param interface{}) {
	//fmt.Println(time.Now().Second())
	//client  := handle.(* BaseClient)
	/*if client.LastHeartBeatCount == 0 {
		client.LastHeartBeatCount = time.Now().Unix()
	}*/
}

// 心跳检查定时器
func HeartBeatCheck(handle interface{}, id int , param interface{}) {
	client  := handle.(* BaseClient)
	curTimeCount := time.Now().Unix()
	if 	client.LastHeartBeatCount != 0 && curTimeCount - client.LastHeartBeatCount > _const.ClientHeartBeatCheckDur {
		// 未检测到心跳断开了 此时应该重连
		if client.ReStart() {
			client.RestartCount ++
			if client.RestartCount >= _const.ClientRestartTolerantTimes { //连续重启3次失败
				fmt.Println("连续3次重启失败.... 应该是服务端 关闭了")
				os.Exit(1)
			}
		}
	}
}

/*
	开启定时器
*/
func (client * BaseClient) CreateTimer(hbdur ,  hbcheck int32) bool {
	//开启定时器 30s心跳检测回调
	client.Timer = &timer2.CronTimer{}
	client.Timer.Add(HeartBeatSend,client, nil , hbdur)
	client.Timer.Add(HeartBeatCheck,client, nil , hbcheck)
	client.Timer.Start()
	return true
}

func (client * BaseClient) Initial() {
	client.CreateTimer(_const.ClientHeartBeatDur, _const.ClientHeartBeatCheckDur)
}


func (client * BaseClient) StopTimer() bool {
	//开启定时器 30s心跳检测回调
	client.Timer.Stop()
	return true
}

func (client * BaseClient) ReStart() bool {
	//连接服务器
	conn,err := net.DialTCP("tcp",nil,client.hawkServer)
	if err != nil {
		fmt.Printf("connect to hawk server error: [%s]",err.Error())
		return false
	}

	client.Conn = conn
	client.stopChan = make(chan struct{})
	client.StartRecv()
	return true
}

func (client * BaseClient) Start(server string) bool{
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
	client.stopChan = make(chan struct{})

	client.StartRecv()
	return true
}

func (client * BaseClient) Send( buf []byte){
	if client.Conn != nil  {
		_ , err := client.Conn.Write(buf)
		if err != nil {
			fmt.Println("send error-->",err.Error())
		}
	}
}

func (client * BaseClient) StartRecv() bool {
	go receivePackets(client)
	return true
}

func (client * BaseClient) OnRecvFrame(buf []byte) {
	fmt.Println("recv string",len(buf[:]),buf[:])
//	client.Send(buf[:])
}

func (client * BaseClient) Stop() {
	client.Conn.Close() //会退出接收go程
	client.StopTimer()
}

func (client * BaseClient) isIpv4() bool {

	return true
}

// 接收数据包
func receivePackets(client *BaseClient) {
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
		client.OnRecvFrame(buf[:n]) //调用回掉
	}
	fmt.Println("Recieve go routinue finished")
}

