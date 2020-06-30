package frontend

import (
	_const "Pushsystem/src/const"
	"Pushsystem/src/pkg/tcpclient"
	"Pushsystem/src/protocol"
	"Pushsystem/src/utils"
)
/*
	设备注册 解析逻辑
*/
func (vServer * VirtualServer) OnGetDeviceRegister(client *tcpclient.TcpClient , protocolData *protocol.Protocol) {
	deviceUniqueID := utils.UniqueId(int32(protocolData.Head.DeviceType) ,string(protocolData.Head.DeviceID[:]))
	//协议中的GateWayID 字段是GateWay服务器计算的值 ，和 virtualServer维护的 tcp连接地址一致
	if cli,ok := vServer.GateWayIDLinksMap.Load(string(protocolData.Head.GateWayID[:])); !ok {
		panic("本地解析和客户端打包逻辑没对应 	protocolData.Head.GateWayID")
	}else {
		if cli != client{
			panic("此时两个实例必须相同")
		}
		//保存 设备唯一标识 和 GateWayID映射关系  todo：需要存储到redis中统一保存
		// todo: 需要本地实现LRU或者LFU 算法保存关系到内存中，减少对redis的访问压力
		vServer.UniqueIDIpMap.Store(deviceUniqueID,protocolData.Head.GateWayID)	
		//返回 注册帧应答帧 //todo：注册镇返回未定义，先简单写一个返回
		vServer.SendRegisterResponse(client,protocolData)
	}
}
/*

*/
func (vServer * VirtualServer) SendRegisterResponse(client *tcpclient.TcpClient , protocolData *protocol.Protocol) {
	_, addr,idc := utils.GetServerInstance()
	copy(protocolData.Head.ManagerID[:] , []byte(utils.MD5(addr))[:])
	protocolData.Head.ManagerIDC = idc
	protocolData.ClearFlag(_const.ProtoFlagUpStream) //设置为下行数据
	//生成唯一包标识ID
	frameID := utils.GetFrameIdFacInstance().GetFrameID()
	copy(protocolData.PackID[:],[]byte(frameID)[:])
	sendBuf := protocolData.Package()
	client.Send(sendBuf)
}



