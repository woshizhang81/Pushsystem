package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)



func main() {
	c, _, err := zk.Connect([]string{"centos-pc1:2181","centos-pc1:2181","centos-pc1:2181"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
	children, stat, ch, err := c.ChildrenW("/")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v %+v\n", children, stat)
	e := <-ch
	fmt.Printf("%+v\n", e)
}



