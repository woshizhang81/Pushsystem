package main

import (
	"github.com/letsfire/factory"
	"fmt"
	"time"
)

func main(){
	var master = factory.NewMaster(20000, 8)

	// 新建第一条工作流水线
	var line1 = master.AddLine(func(args interface{}) {

		// TODO 处理您的业务逻辑
		 fmt.Println("line1",args)
	})

	// 新建第二条工作流水线
	/*var line2 = master.AddLine(func(args interface{}) {

		// TODO 处理您的业务逻辑
		 fmt.Println("line2",args)
	})*/

	master.Running()            // 正在运行的协程工人数量
	go production("line1",line1)
//	go production("line2",line2)
	// 根据业务场景将参数提交

	for{
		number := master.Running()
		fmt.Println("cur number is",number)
		time.Sleep(time.Second)
	}
	/*
	for j := 0; j < 10; j++ {
		line2.Submit(j)
	}
	*/
	// 协程池数量可动态调整
//	master.AdjustSize(100)      // 指定数量进行扩容或缩容
//	master.Shutdown()           // 等于 master.AdjustSize(0)
}

func production(mark string,line *factory.Line){
	var args int64 = 0

	if mark == "line1" {
		args = 100000
	}

	for  {
		line.Submit(args)
		args++
		if args % 100000 == 0 {
			time.Sleep( time.Second)
		}else {
			//time.Sleep( time.Microsecond)
		}
	}
}

