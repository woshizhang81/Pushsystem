package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
)
type songInfo struct {
	Name     string
	Duration int
}

type config struct {
	Bc string
	Song songInfo
}

func test_toml() {
	var cg config
	var cpath string = "/home/work/goworkspace/src/PushSystem/src/example/example.toml"
	if _, err := toml.DecodeFile(cpath, &cg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v %v\n", cg.Bc, cg.Song)
}

func main1() {
	test_toml()
}

type lll struct{
	v  int
	a  int
}
func change (l *lll){
	l.v = 1
	l.a = 2
}
func main()  {
	l := lll{}
	change(&l)
	fmt.Printf("%v\n" ,l)


/*	path := "/home/work/goworkspace/src/PushSystem/conf/gateway.toml"
	config := gateway.GateWayConfig{}
	 _, err := toml.DecodeFile(path,&config)
	 if err != nil{
	 	fmt.Println(err)
	}
	fmt.Printf("%v",config)*/
}
