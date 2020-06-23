package main
import (
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"sync"
	"time"
)
var wg *sync.WaitGroup

func main() {
	conn, _, err := zk.Connect([]string{"centos-pc1:2181","centos-pc2:2181","centos-pc3:2181"}, time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	//zk 包没有提供rmr命令，只能递归删除了
	if b, _, _ := conn.Exists("/demo"); b {
		log.Println("exists /demo")
		paths, _, _ := conn.Children("/demo")
		for _, p := range paths {
			conn.Delete("/demo/"+p, -1)
		}
		err = conn.Delete("/demo", -1)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("delete /demo")
		}
	}
	wg = &sync.WaitGroup{}
	watchDemoNode("/demo", conn)
	wg.Wait()
}
func watchDemoNode(path string, conn *zk.Conn) {
	wg.Add(1)
	//创建
	watchNodeCreated(path, conn)
	//改值
	go watchNodeDataChange(path, conn)
	//子节点变化「增删」
	go watchChildrenChanged(path, conn)
	//删除节点
	watchNodeDeleted(path, conn)
	wg.Done()
}
func watchNodeCreated(path string, conn *zk.Conn) {
	log.Println("watchNodeCreated")
	for {
		_, _, ch, _ := conn.ExistsW(path)
		e := <-ch
		log.Println("ExistsW:", e.Type, "Event:", e)
		if e.Type == zk.EventNodeCreated {
			log.Println("NodeCreated ","path=",e.Path)
			return
		}
	}
}
func watchNodeDeleted(path string, conn *zk.Conn) {
	log.Println("watchNodeDeleted",)
	for {
		_, _, ch, _ := conn.ExistsW(path)
		e := <-ch
		log.Println("ExistsW:", e.Type, "Event:", e)
		if e.Type == zk.EventNodeDeleted {
			log.Println("NodeDeleted ","path=",e.Path)
			return
		}
	}
}
func watchNodeDataChange(path string, conn *zk.Conn) {
	for {
		_, _, ch, _ := conn.GetW(path)
		e := <-ch
		log.Println("GetW('"+path+"'):", e.Type, "Event:", e)
	}
}
func watchChildrenChanged(path string, conn *zk.Conn) {
	for {
		_, _, ch, _ := conn.ChildrenW(path)
		e := <-ch
		log.Println("ChildrenW:", e.Type, "Event:", e)
	}
}