package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	dir,_ := os.Getwd()
	fmt.Println("当前路径：",dir)

//	path := getCurrentPath()
//	fmt.Println(path)
}

func getCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	checkErr(err)
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}