package safemap


import "sync"
type SafeMapBase interface {
 	Get(k string) interface{}
	Set(k string, v interface{}) bool
 	GetSuccessHook(arg interface{}) bool
 	SetSuccessHook(arg interface{}) bool
}

type SafeMap struct {
	Data map[string]interface{}
	Lock sync.RWMutex
	Handle interface{}
}

func (this *SafeMap) GetSuccessHook(arg interface{}) bool{
	return true
}
func (this *SafeMap) SetSuccessHook(arg interface{}) bool{
	return true
}

func (this *SafeMap) Get(k string) interface{} {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	if v, exit := this.Data[k]; exit {
		if true == this.GetSuccessHook(v) {
			return v
		}
	}
	return nil
}

func (this *SafeMap) Set(k string, v interface{}) bool {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	if this.Data == nil {
		this.Data = make(map[string]interface{})
	}
	if true == this.SetSuccessHook(v) {
		this.Data[k] = v
		return true
	}
	return false
}
