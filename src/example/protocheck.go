package main

import (
	"Pushsystem/src/protocol"
	"fmt"
//	"os"
//	"encoding/binary"
	"container/list"
)

type body struct{
	Data [2]byte
}

func ( obj *body)UnPacking(buf []byte) uint16  {
	copy(obj.Data[:],buf[:])
	return uint16(len(obj.Data))
}

func ( obj *body)GetBuffer() []byte   {
	return obj.Data[:]
}//获得动态buffer 返回buffer长度

func PrintByHex(buf []byte){
	for _, v:= range buf {
		fmt.Sprintf("%x",v)
	}
}

func main ()  {

	tbody := body{}
	tbody.Data[0] = 0x80
	tbody.Data[1] = 0x81

	proto := protocol.Protocol{}
	proto.Init()
	proto.Flag[0]   = 0x00
	proto.Flag[1]   = 0x00
	proto.PackType  = 6
	str := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	copy(proto.PackID[:],[]byte(str)[:])

	proto.Head.ModID = 12

	str  = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	copy(proto.Head.ModSerID[:],[]byte(str)[:])
	proto.Head.ModSerIDC = 10

	str  = "cccccccccccccccccccccccccccccccc"
	copy(proto.Head.GateWayID[:],[]byte(str)[:])
	proto.Head.GateWayIDC = 20

	str  = "dddddddddddddddddddddddddddddddd"
	copy(proto.Head.ManagerID[:],[]byte(str)[:])
	proto.Head.ManagerIDC = 30

	str  = "127.0.0.1:13232"
	copy(proto.Head.ClientAddr[:],[]byte(str)[:])

	proto.Data = &tbody

	buf := proto.Package()
	//fmt.Printf("%v", buf)
	fmt.Printf("length = %d ,buf = %v\n", len(buf),buf)

	fmt.Println("00000000000000000000000000000000000000000")
	proto2 := protocol.Protocol{}
	proto2.Data = &tbody
	proto2.UnPacking(buf)
	fmt.Printf("data = %+v\n", proto2)
	fmt.Printf("data = %+v\n", proto2.Data)
	fmt.Println("111111111111111111111111111111111111111111")

	checkObj := &protocol.ProtoCheck{}
	checkObj.Init()
	//var frames [][]byte
	//frames = make([][]byte,0,2)
	list1 := list.New()

	doublebuf := append(buf,buf...)
	checkObj.CheckAndGetProtocolBuffer(doublebuf,list1)
	fmt.Println(list1.Len())

	for item := list1.Front();nil != item ;item = item.Next() {
		//fmt.Println(len(item.Value.([]byte)),item.Value.([]byte))
	//	fmt.Println(item.Value.([]byte))
		//fmt.Printf("%v",item.Value.([]byte))

		newstr := fmt.Sprintf("[% X]",item.Value.([]byte))
		fmt.Println(newstr)
	}

}
