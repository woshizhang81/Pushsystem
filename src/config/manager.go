package config

import (
	"Pushsystem/src/utils"
	"github.com/BurntSushi/toml"
	"log"
	"sync"
)

const ManagerConfigName = "manager.toml"

type ManagerConfig struct{
	Module 	string
	IpVer  	uint8
	Port 	uint16
	IDC		uint16
}


var _managerInstance * ManagerConfig

var onceManager sync.Once

func GetManagerInstance() *ManagerConfig{
	onceManager.Do(func(){
		_managerInstance = &ManagerConfig{}
	})
	return _managerInstance
}

/*
	加载配置
*/
func (obj * ManagerConfig)LoadConfig()  *ManagerConfig {
	filename := utils.GetConfigPath() + ManagerConfigName
	_, err := toml.DecodeFile(filename, obj)
	if err != nil {
		msg := "datadef file :" + filename + "load failed"
		log.Fatal(msg)
		return nil
	}
	return obj
}



