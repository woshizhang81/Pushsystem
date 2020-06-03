package protocol

import "encoding/binary"

//	设备心跳结构体
type HeartBeatBodyDevice struct {

}

func (obj * HeartBeatBodyDevice) GetBuffer() []byte {
	return []byte{}
}

func (obj * HeartBeatBodyDevice) UnPacking(buf [] byte) uint16 {
	return 0
}


//	manager心跳结构体 //携带manager的当前运行参数，用于 计算权重 实现跨机房负载均衡
type HeartBeatBodyManager struct {
	CpuUsage	uint8   	//CPU使用率,百分制度 84 单位%
	MemUsage	uint8   	//内存使用率,百分制度 84 单位%
	UpQps		uint32   		//上行 qps 5000
	DownQps		uint32   		//下行 qps 5000
}

func (obj * HeartBeatBodyManager) GetBuffer() []byte {
	retBuf := make([]byte, 0,6)

	retBuf = append(retBuf , obj.CpuUsage)
	retBuf = append(retBuf , obj.MemUsage)

	binary.BigEndian.PutUint32(retBuf[2:] , obj.UpQps)
	binary.BigEndian.PutUint32(retBuf[4:] , obj.DownQps)
	return retBuf
}


func (obj * HeartBeatBodyManager) UnPacking(buf [] byte) uint16 {
	length := uint16(0)
	obj.CpuUsage  = buf[0]
	obj.CpuUsage  = buf[1]
	length = length + 2
	obj.UpQps = binary.BigEndian.Uint32(buf[length:])
	length = length + 4
	obj.DownQps = binary.BigEndian.Uint32(buf[length:])
	length = length + 4
	return length
}
