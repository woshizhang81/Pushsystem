package config

import (
	"Pushsystem/src/utils"
	"github.com/BurntSushi/toml"
	"log"
	"sync"
)
const AppConfigName = "app.toml"
type appConfig struct{
	AppName string
	Role	string
}

var _instance * appConfig
var once sync.Once

func GetInstance() *appConfig {
	once.Do(func(){
		_instance = &appConfig{}
	})
	return _instance
}

/*
	加载配置
*/
func (obj * appConfig)LoadConfig()  *appConfig{
	filename := utils.GetConfigPath() + AppConfigName
	_, err := toml.DecodeFile(filename,obj)
	if err != nil{
		msg := "datadef file :"+ filename +"load failed"
		log.Fatal(msg)
		return nil
	}
	return obj
}