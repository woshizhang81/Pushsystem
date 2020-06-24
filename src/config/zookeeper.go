package config

import (
	"Pushsystem/src/utils"
	"github.com/BurntSushi/toml"
	"log"
	"sync"
)
const ZkClusterConfigName = "zkcluster.toml"
type ZooKeeperConfig struct{
	Servers []string
}

var _instanceZk * ZooKeeperConfig
var onceZk sync.Once

func GetZkInstance() *ZooKeeperConfig{
	once.Do(func(){
		_instanceZk = &ZooKeeperConfig{}
	})
	return _instanceZk
}

/*
	加载配置
*/
func (obj * ZooKeeperConfig)LoadConfig()  *ZooKeeperConfig{
	filename := utils.GetConfigPath() + ZkClusterConfigName
	_, err := toml.DecodeFile(filename,obj)
	if err != nil{
		msg := "datadef file :"+ filename +"load failed"
		log.Fatal(msg)
		return nil
	}
	return obj
}