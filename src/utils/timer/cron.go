package timer

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Task func (handle interface{} , dur int32,param interface{})

type TimeUnit struct{
	handle     interface{}
	param      interface{}
	task       Task
	dur        int32  //定时间隔
	startCount int64 //当前定时计数
}

type CronTimer struct {
	chanStop chan<-int // 停止信号:= make(chan int)
	chanPerSecond chan<-int //每秒定时信号通道 := make(chan int, 10)
	taskMap sync.Map//线程安全
	running int32    //关闭定时器的标志位 0为运行，非0为停止
}

func (cron *CronTimer) Init() {
	atomic.StoreInt32(&cron.running,0) //设置为0
}

func (cron *CronTimer) Start () {
	atomic.StoreInt32(&cron.running,0) //设置为0
	go run(cron)
}

//添加定时任务,
//params task 回调函数， handle回调函数传入句柄，params 回调函数传入参数 dur定时间隔
func (cron *CronTimer) Add (task Task,handle interface{},params interface{},dur int32) {
	newTask := TimeUnit{}
	newTask.task = task
	newTask.handle = handle
	newTask.startCount = time.Now().Unix()
	newTask.dur = dur
	newTask.param = params
	cron.taskMap.Store(dur,newTask)
}

//
func (cron *CronTimer) Delete(key int32) {
	cron.taskMap.Delete(key)
}

func (cron *CronTimer) Stop () {
	atomic.StoreInt32(&cron.running,1) // 非0即可
}

func execute(Timer * CronTimer , c <-chan int){ //c 用作接收
	//startTimeCount := time.Now().Unix()
	var count int64  //总计数 unit seconds
	for {
		select{
		case sign := <-c:
			if sign == 1 {
				count ++
			    Timer.taskMap.Range(func(key, value interface{}) bool {
			    	realKey := key.(int32)
			    	realValue := value.(TimeUnit)
			    	//startTime := realValue.startCount - startTimeCount
					mod := count % int64(realValue.dur)
					if mod == 0 {
						//到了定时时间了 執行回調
						realValue.task(realValue.handle,realKey,realValue.param)//调用回调
					}
					return true
					})
				} else {
					fmt.Println("channel close")
					break
				}
		case <-time.After(3 * time.Second):
			fmt.Println("time out 3 seconds")
		}
	}
}

func run(timer *CronTimer){
	chanExecute := make(chan int ,10)
	go execute(timer , chanExecute)

	c := time.Tick(time.Second)
	for{
		if timer.running == 0 {
			<-c
			chanExecute <- 1 //每1s的定时信号
		}else{
			close(chanExecute)
			break
		}
	}
}


