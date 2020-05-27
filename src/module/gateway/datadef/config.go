package datadef
type FrontConfig struct {
	QpsLimit int16//qps 限制
	MaxConn int32   //最大客户端连接数字
	SlotNum int16 		//槽位数量
	BlackListNum int16      //黑名单长度
	Ip string				//ip地址
	Port uint16				//端口号
}

type BackConfig struct {
	QpsLimit int16//qps 限制
	MaxConn int32   //最大客户端连接数字
	SlotNum int16 		//槽位数量
	BlackListNum int16      //黑名单长度
	Ip 		string				//ip地址
	Port 	uint16				//端口号
}
/*gateway module 配置*/
type GateWayConfig struct{
	Module string
	Frontend FrontConfig
	Backend	 BackConfig
}

