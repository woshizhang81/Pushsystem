package basenet

import (
	//"Pushsystem/src/pkg/tcpserver"
	"bufio"
	"fmt"
	"net"
	"Pushsystem/src/pkg/tcpserver"
	"strconv"
)

/*
	goroutine per connection
*/
type NetServer struct {
	 callbackHandle  interface{}
	 listener        net.Listener
	 acceptCallback  tcpserver.OnAccept
	 receiveCallback tcpserver.OnReceive
	 closeCallback   tcpserver.OnClose
}

func (server *NetServer) SetCallBackHandle (handle interface{}){
	server.callbackHandle = handle
}

func (server *NetServer)SetAcceptCallback(fun tcpserver.OnAccept){
	server.acceptCallback = fun
}

func (server *NetServer)SetReceiveCallback(fun tcpserver.OnReceive){
	server.receiveCallback = fun
}

func (server *NetServer)SetCloseCallback(fun tcpserver.OnClose){
	server.closeCallback = fun
}

func Construct() *NetServer {
	return &NetServer{}
}

/*
当需要的容量超过原切片容量的两倍时，会使用需要的容量作为新容量。
当原切片长度小于1024时，新切片的容量会直接翻倍。而当原切片的容量大于等于1024时，会反复地增加25%，直到新容量超过所需要的容量。
*/

func process(conn net.Conn, handle *NetServer) {
	defer conn.Close()
	//var data []byte = make([]byte, 4096)   //拼接
	//var buf  []byte = make([]byte, 1305)   //缓存切片
	for {
		var buf  []byte = make([]byte, 1305)   //缓存切片
		//n, err := conn.Read(buf[:])
		reader := bufio.NewReader(conn)
		n , err := reader.Read(buf)
		if n == 0 {
			handle.closeCallback(handle.callbackHandle,conn)
			break
		}
		if err != nil {
			handle.closeCallback(handle.callbackHandle,conn)
			fmt.Printf("read from connect failed, err: %v\n", err)
			break
		}
		//data = append(data,buf...)  //... 切片打散
		handle.receiveCallback(handle.callbackHandle,conn,buf[:n-1])
		//		fmt.Printf("receive from client, data: %v\n", buf)
	}
}

func (server *NetServer ) Create (ipAddr string, port uint16) bool {
	address := ipAddr + ":" + strconv.Itoa(int(port))
	fmt.Println(address)
	listener, err := net.Listen("tcp", address)
	server.listener = listener
	if err != nil {
		fmt.Printf("listen fail, err: %v\n", err)
		return false
	}
	//2.accept client request
	//3.create goroutine for each request
	for {
		conn, err := server.listener.Accept()
		if err != nil { //此处要判断的还有信号打断的情况
			fmt.Printf("accept fail, err: %v\n", err)
			break
		}
		//create goroutine for each connect
		go process(conn, server)
		server.acceptCallback(server.callbackHandle,conn)
	}
	return true
}

func (server *NetServer)Send (conn net.Conn, buf []byte){

}

func (server *NetServer) ShutDown(){
	server.listener.Close()
}

