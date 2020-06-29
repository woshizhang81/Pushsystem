package zkclient

import (
	_const "Pushsystem/src/const"
	"Pushsystem/src/utils"
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"sync"
	"time"
)
/*
	zk 客户端的开启
	zk 节点操作的封装
	zk 节点所有事件的监听封装
	回调方式展示，每个节点的每个事件对应的改变信息
*/
const(
	ZKChildAdd = 1
	ZKChildDel = 0
)

type OnPathCreated func (handle interface{},path string ,nodeValue []byte)                   //该节点被创建或者删除
type OnPathDeleted func (handle interface{},path string)
type OnPathContextChanged func (handle interface{}  , path string, latestPathValue []byte ,currentPathValue []byte)
type OnPathChildNumChanged func (handle interface{} , path string, changeType uint8 , ChangedNode string) //

type ZkNodeEvents  struct {
	wg *sync.WaitGroup
	NodeName string 	//节点名称
	Conn	*zk.Conn
	OldNodeEventContext NodeEventContext    //最近一次的 该path的值的记录。采用值传递方式 , 次结构要作为sync.map value
	CallBackHandle interface{} //回调句柄
	CallBackPathCreated OnPathCreated
	CallBackPathDeleted OnPathDeleted
	CallBackValueChanged OnPathContextChanged
	CallBackChildChanged OnPathChildNumChanged
	EventsStopChan [_const.ZookeeperEventNumber]chan bool
 }

func NewZkNodeEvents(nodeName string, conn *zk.Conn, callBackHandle interface{}, callBackPathCreated OnPathCreated, callBackPathDeleted OnPathDeleted, callBackValueChanged OnPathContextChanged, callBackChildChanged OnPathChildNumChanged) *ZkNodeEvents {
	ret := &ZkNodeEvents{NodeName: nodeName, Conn: conn, CallBackHandle: callBackHandle, CallBackPathCreated: callBackPathCreated, CallBackPathDeleted: callBackPathDeleted, CallBackValueChanged: callBackValueChanged, CallBackChildChanged: callBackChildChanged}
	for i:=0;i<_const.ZookeeperEventNumber ; i++  {
		ret.EventsStopChan[i] = make(chan bool)
	}
	return ret
}

func watchNodeExist(events * ZkNodeEvents, Id uint8) {
	log.Println("watchNodeCreated")
	for {
		_, _, ch, err := events.Conn.ExistsW(events.NodeName)
		if err != nil {
			log.Println("")
			break
		}
		select {
		case _, ok := <-events.EventsStopChan[Id]:
			if !ok {
				break
			}
		case e :=<-ch :
			log.Println("ExistsW:", e.Type, "Event:", e)
			value,_,err := events.Conn.Get(e.Path)
			if err != nil {
				panic("this cannot happend")
			}
			if e.Type == zk.EventNodeCreated {
				events.CallBackPathCreated(events.CallBackHandle,e.Path,value)
				log.Println("NodeCreatedCreated  ","path=",e.Path)
			} else if e.Type == zk.EventNodeDeleted{
				events.CallBackPathDeleted(events.CallBackHandle,e.Path)
				log.Println("NodeCreatedCreated  ","path=",e.Path)
			}else {
				panic("this cannot happend ..")
				break
			}
		}
	}
	events.wg.Done()
}

func watchNodeDataChange(events * ZkNodeEvents,Id uint8) {
	for {
		_, _, ch, err:= events.Conn.GetW(events.NodeName)
		if err != nil {
			log.Println("")
			break
		}
		select {
		case _, ok := <-events.EventsStopChan[Id]:
			if !ok {
				break
			}
		case e :=<-ch :
			value,_,err := events.Conn.Get(e.Path)
			if err != nil {
				panic("this cannot happend")
			}
			events.CallBackValueChanged(events.CallBackHandle,e.Path,events.OldNodeEventContext.NodeContext,value)
			log.Println("GetW('"+events.NodeName+"'):", e.Type, "Event:", e)
		}
	}
	events.wg.Done()
}

func watchChildrenChanged(events * ZkNodeEvents ,Id uint8) {
	for {
		_, _, ch, err:= events.Conn.ChildrenW(events.NodeName)
		if err != nil {
			log.Println("")
			break
		}
		select {
		case _, ok := <-events.EventsStopChan[Id]:
			if !ok {
				break
			}
		case e :=<-ch :
			CurChildes,_,err := events.Conn.Children(e.Path)
			if err != nil {
				panic("this cannot happend")
			}
			OldChildes := events.OldNodeEventContext.ChildNodeNames
			diffChildArray  := utils.FindDifferentSlice(OldChildes,CurChildes)
			var eventType uint8
			if len(CurChildes) > len (OldChildes) {
				//比原来增加了 找出
				eventType = ZKChildAdd
			}else if len(CurChildes) < len(OldChildes) {
				//比软来减少了，找出减少的 childPath
				eventType = ZKChildDel
			}else {
				panic("this cannot happend")
			}
			for i := 0;i< len(CurChildes) ; i++  {
				events.CallBackChildChanged(events.CallBackHandle,e.Path,eventType,diffChildArray[i])
			}
			log.Println("ChildrenW:", e.Type, "Event:", e)
		}
	}
	events.wg.Done()
}


func (events * ZkNodeEvents)GetCurNodeOldContext(){
	//先把对应的节点数据拿到
	//todo: 每个机房 绑定events.OldNodeEventContext.IDCNodeMap = make(map[uint16] *NodeObject)
	events.OldNodeEventContext.NodeName = events.NodeName
	//获取该节点内容
	value,_,err := events.Conn.Get(events.NodeName)
	if err != nil {
		panic("this cannot happend")
	}
	events.OldNodeEventContext.NodeContext= value
	// 获取该节点的child列表
	childes,_, errCh := events.Conn.Children(events.NodeName)
	if errCh != nil {
		panic("this cannot happend")
	}
	events.OldNodeEventContext.ChildNodeNames = childes
}

func (events * ZkNodeEvents)Run(){
	events.GetCurNodeOldContext()
	//创建监听事件回掉
	events.wg = &sync.WaitGroup{}
	go func(conn *zk.Conn) {
		events.wg.Add(_const.ZookeeperEventNumber)
		go watchNodeExist(events,1)
		go watchNodeDataChange(events,2)
		go watchChildrenChanged(events ,3)
		events.wg.Wait()
		log.Println("path:",events.NodeName, "Monitor Finished")
	}(events.Conn)
}
// 节点通用json结构定义 todo： 根据具体业务补充 目前定义如下 两个域
type NodeObject struct{
	Idc uint16 	//IDC 机房代号
	IdcName	string // 机房字符串 ”sz“ ”gz“  应该与Idc一一对应
}

type NodeEventContext struct {
	NodeContext    []byte   //节点内容  todo:我们约定内容为json 字符串 目前约定{idc="sz"}
	ChildNodeNames []string //子结点名称
	//IDCNodeMap	   map[uint16]*NodeObject
	NodeName       string
}

/*
	return
*/
func (nodeEventObj *NodeEventContext)GetIdcInfo() *NodeObject {
	nodeObj := &NodeObject{}
	err := json.Unmarshal(nodeEventObj.NodeContext,nodeObj)
	if err != nil {
		fmt.Println("json parser failed")
		return nil
	}
	return nodeObj
}

func NewNodeEventContext(nodeName string) *NodeEventContext {
	retHandle := &NodeEventContext{NodeName: nodeName}
	retHandle.NodeContext = make([]byte,0,_const.ZookeeperNodeValueSize) //固定该ZK 所有节点大小最大为1024
	retHandle.ChildNodeNames = make([]string,0,_const.ZookeeperChildNodeSize) //固定该ZK 所有节点大小最大为1024
	return retHandle
}

type ZkClient struct {
	WorkState bool //工作状态
	ZkAddr []string  // zookeeper 集群地址
	Conn	*zk.Conn
	PathEventsMap sync.Map //每条路径的事件映射关系 path ： ZkEvents
}

func (client * ZkClient) Initial(){
	//loadConfig //获得zookeeper 地址列表
}

func (client * ZkClient) Start()  bool {
	conn, _, err := zk.Connect(client.ZkAddr , time.Second)
	client.Conn = conn
	if err != nil {
		return false
	}
	return true
}

/*
	添加该路径所有事件
*/
func (client * ZkClient) AddPathEvents(
	Path string,
	Handle 	interface{},
	CallBackPathCreated OnPathCreated,
	CallBackPathDeleted OnPathDeleted,
	CallBackPathContextChanged OnPathContextChanged,
	CallBackPathChildNumChanged OnPathChildNumChanged) {
	newNode := NewZkNodeEvents(Path,client.Conn, Handle,
				CallBackPathCreated, CallBackPathDeleted,
				CallBackPathContextChanged, CallBackPathChildNumChanged)
	newNode.Run()
	client.PathEventsMap.Store(Path,newNode)
}

/*
   删除该路径的所有事件
*/
func (client *ZkClient) DeletePathEvents(path string){
	value, ok := client.PathEventsMap.Load(path)
	if ok {
		val := value.(*ZkNodeEvents)
		for i := 0 ; i < _const.ZookeeperEventNumber ; i++  {
			close(val.EventsStopChan[i])
		}
	}
	client.PathEventsMap.Delete(path)
}
/*
	停止该zkclient
*/
func (client *ZkClient) Stop() {
	client.Conn.Close()
	client.PathEventsMap.Range(func(key, value interface{}) bool {
		path := key.(string)
		client.DeletePathEvents(path)
		client.PathEventsMap.Delete(path)
		return true
	})
}
