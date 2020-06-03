package main

import "fmt"

type IIII interface {
	toString()
}

type Test struct {

}

func (obj * Test) toString(){
	fmt.Println("kjljkjljljljl")
}

func main(){

	var temp IIII
	temp = &Test{}
	temp.toString()
}
