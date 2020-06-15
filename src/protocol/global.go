package protocol

import (
	"Pushsystem/src/const"
	"encoding/binary"
	"Pushsystem/src/utils"
)

/*
	frame 必须是经过校验的帧
	gateway 直接从帧里 按偏移 读出 类型 ，提高效率
*/
func GetTypeFromFrame(frame []byte) uint8 {
	//偏移为
	PackType := frame[_const.GateWayProtocalTypeOffset]
	return PackType
}


/*
	frame 必须是经过校验的帧
	gateway 直接从帧里 按偏移 读出 DeviceID DeviceType，提高效率
*/
func GetUniqueIdFromFrame(frame []byte) string {
	//偏移为
	DeviceID 	:= 	string(frame[_const.GateWayProtocalDeviceIDOffset:50])
	DeviceType 	:= 	binary.BigEndian.Uint16(frame[_const.GateWayProtocalDeviceTypeOffset:])
	UniqueID := utils.UniqueId(int32(DeviceType),DeviceID)
	return UniqueID
}
