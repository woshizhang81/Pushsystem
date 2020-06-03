package protocol

import "encoding/binary"


//md5(deviceid + type + pagname + timestamp)
type RegisterBodyDevice struct {
	PkgName		[255]byte	//包名 , 最大255字节, 不够补0
	TimeStamp	int64		//时间戳 注册的
	Token 		[50]byte	//Token  //需要实现 TOKEN 服务和check机制 最大50字节
}

func (obj * RegisterBodyDevice) GetBuffer() []byte {
	retBuf := make([]byte,0,313) //255 + 8 +50
	retBuf = append(retBuf , obj.PkgName[:]...)
	binary.BigEndian.PutUint64(retBuf[255:] , uint64(obj.TimeStamp))
	retBuf = append(retBuf , obj.Token[:]...)
	return retBuf
}

func (obj * RegisterBodyDevice) UnPacking(buf [] byte) uint16 {
	length := uint16(0)
	copy(obj.PkgName[:],buf[:255])
	length = length + 255
	obj.TimeStamp = int64(binary.BigEndian.Uint64(buf[length:]))
	length = length + 8
	copy(obj.Token[:],buf[length:])
	return length
}


type RegisterBodyManager struct {
	TimeStamp	int64		//时间戳 注册的
	Sign [32]byte	//Token md5(idc + managerId +  TimeStamp)
}

func (obj * RegisterBodyManager) GetBuffer() []byte {
	retBuf := make([]byte,0,40) //
	binary.BigEndian.PutUint64(retBuf[:] , uint64(obj.TimeStamp))
	retBuf = append(retBuf , obj.Sign[:]...)
	return retBuf
}

func (obj * RegisterBodyManager) UnPacking(buf [] byte) uint16 {
	length := uint16(0)
	obj.TimeStamp = int64(binary.BigEndian.Uint64(buf[length:]))
	length = length + 8
	copy(obj.Sign[:],buf[length:])
	return length
}



