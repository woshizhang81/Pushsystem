package utils

import (
	"sync"
	"sync/atomic"
)

/*
	获得本地 网络地址信息 NET ID , 和IDC机房信息
	从配置文件中读取
*/
var _FrameInstance * FrameIDFac

var onceFrameID sync.Once

type FrameIDFac struct{
	uniqueIndex uint64
}

func (frameID *FrameIDFac) GetFrameID() string {
	_,addrStr,_ := GetServerInstance()
	newNum := atomic.AddUint64(&frameID.uniqueIndex , 1)
	uniqueString := addrStr+ string(newNum)
	return MD5(uniqueString)
}

func GetFrameIdFacInstance() *FrameIDFac{
	onceFrameID.Do(func() {
		_FrameInstance = &FrameIDFac{}
	})
	return _FrameInstance
}


