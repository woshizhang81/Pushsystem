package protocol

import (
	"Pushsystem/src/utils"
	"encoding/binary"
	"os"
	"fmt"
	"Pushsystem/src/const"
)
const BufDefaultSize = 2048

type BodyData interface { //统一为PB格式
	GetBuffer() []byte   //获得动态buffer 返回buffer长度
	UnPacking( buf []byte) uint16  //解包
}

/*协议固定场景经过的 节点信息*/
type TransHead struct {
	DeviceID   	[_const.DeviceIDSize]byte //客户端移动端的唯一ID
	DeviceType  uint8  	 //为兼容多种终端唯一标识
	ModID     	uint16   //支持最大65535 个应用
	ModSerID  	[_const.CommonServerIDSize]byte //网关服务器地址  md5 字符串 兼容ip6  v4(v6) md5([ip:port])
	ModSerIDC  	uint16   //网关代表机房
	GateWayID  	[_const.CommonServerIDSize]byte //网关服务器地址
	GateWayIDC  uint16   //网关代表机房
	ManagerID  	[_const.CommonServerIDSize]byte //解析服务器地址
	ManagerIDC  uint16   //解析服务器机房
	ClientAddr 	[_const.NetNodeAddrSize] byte //客户终端端的ip 兼容ip6  v4(v6) [ip:port] 模式
}

func (transHead * TransHead) UnPackage(buf []byte) {
	index := uint16(0)
	copy(transHead.DeviceID[:] ,buf[index:])//[32]byte //网关服务器地址  md5 字符串 兼容ip6 v4(v6) md5([ip:port])
	index = index + 50
	transHead.DeviceType  = buf[index]  //uint8  	 //为兼容多种终端唯一标识
	index ++
	transHead.ModID     	= binary.BigEndian.Uint16(buf[index:])//uint16   //支持最大65535 个应用
	index = index + 2
	copy(transHead.ModSerID[:] ,buf[index:])//[32]byte //网关服务器地址  md5 字符串 兼容ip6 v4(v6) md5([ip:port])
	index = index + 32
	transHead.ModSerIDC     = binary.BigEndian.Uint16(buf[index:])//uint16   //网关代表机房
	index = index + 2
	copy(transHead.GateWayID[:] ,buf[index:])//[32]byte //网关服务器地址  md5 字符串 兼容ip6 v4(v6) md5([ip:port])
	index = index + 32
	transHead.GateWayIDC	= binary.BigEndian.Uint16(buf[index:])//uint16   //网关代表机房
	index = index + 2
	copy(transHead.ManagerID[:] ,buf[index:])//[32]byte //网关服务器地址  md5 字符串 兼容ip6 v4(v6) md5([ip:port])
	index = index + 32
	transHead.ManagerIDC	= binary.BigEndian.Uint16(buf[index:])//uint16   //网关代表机房
	index = index + 2

	copy(transHead.ClientAddr[:] ,buf[index:])//[32]byte //网关服务器地址  md5 字符串 兼容ip6 v4(v6) md5([ip:port])
	index = index +50
	if index != 205 {
		os.Exit(1)
		fmt.Println("head unpack error !!")
	}
}
/*
	头部打包
*/
func (transHead * TransHead) Package() []byte{
	var buf [205]byte  //固定长度205 头部

	tlen := uint8(0)
	copy(buf[tlen:] ,transHead.DeviceID[:])
	tlen = tlen + 50
	buf[tlen] = transHead.DeviceType
	tlen ++

	binary.BigEndian.PutUint16(buf[tlen:] , transHead.ModID)
	tlen = tlen + 2
	copy(buf[tlen:] ,transHead.ModSerID[:])
	tlen = tlen + 32
	binary.BigEndian.PutUint16(buf[tlen:],transHead.ModSerIDC)
	tlen = tlen + 2

	copy(buf[tlen:] ,transHead.GateWayID[:])
	tlen = tlen + 32
	binary.BigEndian.PutUint16(buf[tlen:],transHead.GateWayIDC)
	tlen = tlen + 2
	copy(buf[tlen:] ,transHead.ManagerID[:])
	tlen = tlen + 32
	binary.BigEndian.PutUint16(buf[tlen:],transHead.ManagerIDC)
	tlen = tlen + 2

	copy(buf[tlen:] , transHead.ClientAddr[:])
	tlen = tlen +50
	return buf[:]
}

type Protocol struct {   //传输信息 标准格式定义
	PackHead  	[2]byte  //固定 0xFF 0xFF
	PackSize    uint16   //包大小
	Flag        [2]byte  //标志位 0x0001 上下行数据 1 为上行 0 为下行
						 //0x0002 是否需要应答 1.需要应答 0.不需要
						 //0x0004 是否需要geteway拆包 //用于注册帧和心跳帧 gateway状态同步的  0x0005 生效
						 //0x0008
	PackType  	uint8    //预留255种类型  注册帧，心跳帧，转发,推送  //相对manager服务而言
	PackID		[32]byte // 该包的唯一标识。作推送确认使用
	Head		TransHead //
	Data     	BodyData //有效数据
	Crc         uint16//所有数据做校验
}

/*
	初始化
*/
func (proto *Protocol) Init(){
	proto.PackHead[0] = 0xff
	proto.PackHead[1] = 0xff
}

/*
 打包函数
*/
func (proto *Protocol) Package() []byte  {
	buf := make([]byte , 0, BufDefaultSize) //创建2048 的切片
	buf = append(buf, proto.PackHead[:]...)
	buf = append(buf, 0,	0)
	packSizeIndex := len(buf) -2
	buf = append(buf, proto.Flag[:]...)
	buf = append(buf,proto.PackType)
	buf = append(buf,proto.PackID[:]...)

	head := proto.Head.Package()
	buf = append(buf,head[:]...)

	dataBuf := proto.Data.GetBuffer()
	buf = append(buf,dataBuf[:]...)

	proto.PackSize = uint16(len(buf) + 2)
	binary.BigEndian.PutUint16(buf[packSizeIndex:] , proto.PackSize)

	proto.Crc = utils.Crc16(buf[:])
	buf = append(buf,0,0)
	binary.BigEndian.PutUint16(buf[len(buf) - 2 :] , proto.Crc)
	//填充PackSize

	return buf
}



/*
	是否为上行数据 ?
*/
func (transHead * Protocol) IsUpstream()bool {
	ret := transHead.Flag[1] & 0x01
	if ret == 1 {
		return true
	}else {
		return false
	}
}

/*
	是否需要应答
*/
func (transHead * Protocol) NeedAck() bool {
	ret := transHead.Flag[1] & 0x02
	if ret == 1 {
		return true
	}else {
		return false
	}
}
/*
 	解包函数 接收到的信息,解包
	buf 是已经通过校验的
*/
func (proto *Protocol) UnPacking(buf []byte) () {
	length := uint16(len(buf))
	index := uint16(0)
	copy(proto.PackHead[:],buf[index:])
	index = index + 2
	proto.PackSize = binary.BigEndian.Uint16(buf[index:])
	if length != proto.PackSize {
		fmt.Println("unpacking failed")
		os.Exit(1)
	}
	index = index + 2
	copy(proto.Flag[:],buf[index:])
	index = index + 2
	proto.PackType = buf[index]
	index ++
	copy(proto.PackID[:],buf[index:])
	index = index + 32

	proto.Head.UnPackage(buf[index:])
	index = index + 205

	lengthdata := proto.Data.UnPacking(buf[index:])
	index = index + lengthdata
	proto.Crc = binary.BigEndian.Uint16(buf[index:])
}




