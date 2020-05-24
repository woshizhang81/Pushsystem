package main

import (
	"fmt"
	"reflect"
	"time"
)

func process(c1,c2,c3 chan<- int)  {
	for {
		time.Sleep(time.Second)
		c1 <- 1
		time.Sleep(time.Second)
		c1 <- 1
		time.Sleep(time.Second)
		c2 <- 2
		time.Sleep(3 * time.Second)
		c3 <- 3
	}

}

func TestSelect() {
	c1 := make(chan int)
	c2 := make(chan int, 10)
	c3 := make(chan int, 20)

	go process(c1,c2,c3)
	/*
	go func(c1, c2, c3 chan<- int) {
		for {

			time.Sleep(1 * time.Second)
			c1 <- 1
			time.Sleep(1 * time.Second)
			c1 <- 1
			time.Sleep(1 * time.Second)
			c1 <- 2
			time.Sleep(1 * time.Second)
			c1 <- 3
		}

	}(c1, c2, c3)
	*/
	for {
		select {
		case int1 := <-c1:
			fmt.Println("c1 value :", int1)
		case int2 := <-c2:
			fmt.Println("c2 value :", int2)
		case int3 := <-c3:
			fmt.Println("c3 vaule :", int3)
		case <-time.After(2 * time.Second):
			fmt.Println("timeount")
		}

	}
}

func main(){
	timeCount := time.Now().Unix()
	reflect.TypeOf(timeCount)
	fmt.Println(timeCount,"--",reflect.TypeOf(timeCount))
	//TestSelect()
}