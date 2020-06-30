package utils

import (
	"Pushsystem/src/config"
	"net"
)

/*
	获得本地 网络地址信息 NET ID , 和IDC机房信息
	从配置文件中读取
*/
func GetServerInstance() ( ok bool ,ipAddr string, idcId uint16) {
	managerConf := config.GetManagerInstance()
	obj := managerConf.LoadConfig()
	if obj != nil {
		localIp := getLocalIP(obj.IpVer)
		if localIp != ""{
			return true, "localIp" + ":"+ string(obj.Port) , obj.IDC
		}
	}
	return false,"",0
}
// 获取本机网卡IP
// ipType 4=ipv4 6=ipv6
func getLocalIP(ipType uint8) string {
	var (
		err 	error
		adders  []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIpNet bool
	)
	// 获取所有网卡
	 adders, err = net.InterfaceAddrs()
	if  err != nil {
		return  ""
	}
	// 取第一个非lo的网卡IP
	for _, addr = range adders {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				if ipType == 4 {
					return ipNet.IP.String() // 192.168.1.1
				}
			}else if ipNet.IP.To16() != nil {
				if ipType == 6 {
					return ipNet.IP.String()
				}
			}
		}
	}
	return  ""
}


