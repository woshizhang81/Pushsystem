package backend

import (
	"sync"
	"Pushsystem/src/utils"
)

const DefaultManagerNums = 3000  //假设3000个manager

type SafeMap struct {
	AverageIndex uint32	  //平均算法的计数器
	KeysArray 	[]string  //利用切片储所有的key值 ，方便轮训发送或者 计算权重
						  // 由于终端不多，不会影响效率
						  // 保证连续存储key值
	Data map[string]interface{}
	Lock sync.RWMutex
}

func (this *SafeMap)Init(){
	this.KeysArray = make([]string,0,DefaultManagerNums)
}

func (this *SafeMap) GetSize()  int {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	return len(this.KeysArray)
}

/*
	机会均等的获取 map的 value
	todo: 需要根据机房判断，固定idc 一定qps下 平均发送。
	todo: 如果该机房负荷高了,则需要计算权重返回 value 可能是对应的别的机房 （具体算法未实现）
*/
func (this *SafeMap) GetAverage() interface{} {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	this.AverageIndex ++

	if this.AverageIndex == 0xFFFFFFFF {
		this.AverageIndex = 0
	}
	curIndex := this.AverageIndex % uint32(len(this.KeysArray))
	return this.Data[this.KeysArray[curIndex]]
}

func (this *SafeMap) Get(k string) interface{} {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	if v, exit := this.Data[k]; exit {
		return v
	}
	return nil
}

func (this *SafeMap) Set(k string, v interface{})  bool {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	if this.Data == nil {
		this.Data = make(map[string]interface{})
	}
	if _, exit := this.Data[k]; !exit {  //如果不存在才添加
		this.Data[k] = v
		this.KeysArray = append(this.KeysArray, k)
		return true
	}
	return false

}


func (this *SafeMap) Delete(k string)  bool {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	if _, exit := this.Data[k]; exit {  //如果存在才删除
		delete(this.Data,k)
		ret := utils.DeleteValueFormSlice(this.KeysArray,k)
		return ret
	}
	return false
}

type Rangecallback func (key,value interface{}) bool
func (this *SafeMap) Range(cbfun Rangecallback) {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	for key, val := range this.Data{
		if false == cbfun(key, val) {
			break
		}
	}
}

