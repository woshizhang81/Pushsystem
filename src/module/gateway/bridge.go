package gateway

/*
	负责前后端,的数据透传
*/
const UpGoRoutineNum = 100   //上行go程数
const DownGoRoutineNum = 100 //下行go程数
const ChanDepth = 10000   		//每个go程 开辟 Chan的大小
//前端接收到数据,通过channel发送到后端
type BridgeManager struct{


}
