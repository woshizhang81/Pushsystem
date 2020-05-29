package main

import (
	"Pushsystem/src/protocol"
	"Pushsystem/src/utils"
	"fmt"
	"os"
)

type body1 struct{
	Data [2]byte
}
func ( obj *body1)UnPacking(buf []byte) uint16  {
	copy(obj.Data[:],buf[:])
	return uint16(len(obj.Data))
}

func ( obj *body1)GetBuffer() []byte   {
	return obj.Data[:]
}

func main(){
	var lll []byte
	fmt.Println(lll)
	os.Exit(1)
	tbody := body1{}
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

	newbuf := []byte{0,1,2,3}
	newbuf = append(newbuf,buf...)

	fmt.Println(newbuf)
	ret := utils.FindSubByteArray(newbuf,[]byte{0xff ,0xff })
	fmt.Println(ret)
}
