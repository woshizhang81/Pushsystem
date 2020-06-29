package _const

const CommonServerIDSize = 32  //通用服务ID长度
const NetNodeAddrSize = 50     //网络节点 网络地址 长度     ip4(6):port
const DeviceIDSize = 50        //设备ID的长度

const RecieveBufLimit = 1305
const GateWaySlotNum = 500   // slot 个数
const GateWayFrontHbDur = 10 // 前端心跳间隔
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

const UpStreamGoRoutineMaxNum 		= 20000  //最大20000个go程池
const UpStreamGoRoutineInitNum 		= 50	   //初始化50个go程大小
const UpStreamMessageQueueCapity 	= 100000	   //十万个大小队列




const ClientRestartTolerantTimes = 3  //重启容错次数
const ClientHeartBeatCheckDur 	=    GateWayFrontHbDur  // 与端心跳间隔相同
const ClientHeartBeatDur 	  	= 	 GateWayFrontHbDur - 5  // 比前端心跳间隔少5秒，考虑延迟情况


const ZookeeperEventNumber = 3 		//zookeeper 每个路径可监控的事件为3个
const ZookeeperNodeValueSize = 4*4096 		//zookeeper 每个节点最大数据量为4k
const ZookeeperChildNodeSize = 100			//zookeeper 每个节点最大子结点数量

const (
	ZkServerRootName  				= "/push_system"		//zk 根节点名称
	ZkGateWayParentNodeName			= "/push_system/server_gateway"  //zk gateway 父节点名
	ZkManagerParentNodeName			= "/push_system/server_manager" //zk manager 父节点名
)
