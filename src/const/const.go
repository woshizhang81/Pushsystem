package _const

const GateWaySlotNum = 500   // slot 个数
const GateWayFrontHbDur = 30 // 前端心跳间隔
const GateWayFrontFlowRateDur = 5 // 前端流量检测间隔 5s
const GateWayBackHbDur = 30 // 后端心跳间隔
const GateWayBackLoadBalanceDur = 30 // 前端心跳间隔
const GateWayHostStateCheckDur = 10 // 10秒钟检查一次 主机状态

const GateWayProtocalTypeOffset = 6 // 第七字节，偏移为6
const GateWayProtocalDeviceIDOffset 	= 	39 	// 第40字节，偏移为39
const GateWayProtocalDeviceTypeOffset 	= 	89 	// 第90字节，偏移为89



const DownStreamGoRoutineMaxNum 	= 2000  //最大20000个go程池
const DownStreamGoRoutineInitNum 	= 50	   //初始化50个go程大小
const DownMessageQueueCapity 		= 100000	   //十万个大小队列
