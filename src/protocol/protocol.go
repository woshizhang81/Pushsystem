package protocol

import (
	"Pushsystem/src/utils"
	"encoding/binary"
)
const BufDefaultSize = 2048

type BodyData interface {
	GetBuffer() ([]byte ,uint16)  //获得动态buffer 返回buffer长度
}

/*协议固定场景经过的 节点信息*/
type TransHead struct {
	ModID     	uint16   //支持最大65535 个应用
	ModSerID  	[32]byte //网关服务器地址  md5 字符串 兼容ip6  v4(v6) md5([ip:port])
	ModSerIDC  	uint16   //网关代表机房
	GateWayID  	[32]byte //网关服务器地址
	GateWayIDC  uint16   //网关代表机房
	ManagerID  	[32]byte //解析服务器地址
	ManagerIDC  uint16   //解析服务器机房
	DeviceID   	[50]byte //客户端移动端的唯一ID
	DeviceType  uint8  	 //为兼容多种终端唯一标识
	ClientAddr 	[50]byte //客户终端端的ip 兼容ip6  v4(v6) [ip:port] 模式
}
/*
	头部打包
*/
func (transHead * TransHead) Package(buf []byte) uint8{
	len := uint8(0)
	binary.BigEndian.PutUint16(buf[len:2] , transHead.ModID)
	len = len + 2
	copy(buf[len:32] ,transHead.ModSerID[:])
	len = len + 32
	binary.BigEndian.PutUint16(buf[len:2],transHead.ModSerIDC)
	len = len + 2
	copy(buf[len:32] ,transHead.GateWayID[:])
	len = len + 32
	binary.BigEndian.PutUint16(buf[len:2],transHead.GateWayIDC)
	len = len + 2
	copy(buf[len:32] ,transHead.ManagerID[:])
	len = len + 32
	binary.BigEndian.PutUint16(buf[len:2],transHead.ManagerIDC)
	len = len + 2
	copy(buf[len:50] ,transHead.DeviceID[:])
	len = len + 50
	buf[len] = transHead.DeviceType
	len ++
	copy(buf[len:50] , transHead.ClientAddr[:])
	len = len +50

	return len
}

type Protocol struct {   //传输信息 标准格式定义
	PackHead  	[2]byte  //固定 0xFF 0xFF
	PackSize    uint16   //包大小
	Flag        [2]byte  //标志位 0x0001 上下行数据 1 为上行 0 为下行
						 //0x0002 是否需要应答 1.需要应答 0.不需要
						 //0x0004 是否需要geteway拆包 //用于注册帧和心跳帧 gateway状态同步的  0x0005 生效
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
func (proto *Protocol) Package() ([]byte , uint16) {
	buf := make([]byte , BufDefaultSize, BufDefaultSize) //创建2048 的切片
	len := uint16(0)
	copy(buf[len:2] , proto.PackHead[:])
	len = len + 2
	packSizeIndex := len
	len = len + 2

	copy(buf[len:2] , proto.Flag[:])
	len = len + 2

	buf[len] = proto.PackType
	len ++

	copy(buf[len:32] , proto.PackID[:])
	len = len + 32

	headlen := proto.Head.Package(buf[len:])
	len = len + uint16(headlen)

	dataBuf,dataLen := proto.Data.GetBuffer()
	copy(buf[len:],dataBuf[:])
	len = len + dataLen

	proto.Crc = utils.Crc16(buf[:len])
	binary.BigEndian.PutUint16(buf[len:2] , proto.Crc)
	len = len  + 2
	//填充PackSize
	proto.PackSize = len
	binary.BigEndian.PutUint16(buf[packSizeIndex:2] , proto.PackSize)
	return buf,len
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

}




