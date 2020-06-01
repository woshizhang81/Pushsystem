package main
import "icode.baidu.com/baidu/gdp/gdp/logger"
func main(){
	log1 := logger.GetWriter("zhangjj")
	log1.Info(">>>>>>>>>>>>>>>>>>>>>>>")
	log := logger.GetWriter("testlog")

	log.Debug("test")
	log.Trace("test")

	// 以下代码均会在./testlog.log文件中进行输出1行
	log.Info("test")
	log.Notice("test")

	// 以下代码均会在./testlog.log.wf文件中进行输出1行
	log.Warn("test")
	log.Error("test")
	log.Critical("test")

	// 避免日志写完前进程退出
	log.Close()
	log1.Close()

}
