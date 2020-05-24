package main

import (
	"PushSystem/src/utils"
	"fmt"
	"os"
)

func printSlice ( buf []byte) {
	len := len(buf)
	for i:= 0 ;i < len ; i++  {
		fmt.Println(buf[i])
	}
	fmt.Println("\n")
}

func main1(){
//v1 := []int{1, 2} //len=2, cap=2
v1 := make([]byte,2)
v1[0] = 1
v1[1] = 2
fmt.Println(v1)
os.Exit(1)

v1 = append(v1,6,7) //增加一个
v1 = append(v1,6,7,8) //增加多个
fmt.Println(v1,len(v1), cap(v1))  //[1 2 6 7 6 7 8] len=7,cap=8 cap按照初始化的cap倍数增加


//删除一个元素
v2 := []int{1, 2, 3 ,4 ,5, 6}
fmt.Println(v2,len(v2), cap(v2)) //[1 2 3 4 5 6] len=6 cap=6

copyv := append(v2[:1],v2[3:]...) //得到删除后的切片
fmt.Println(copyv,len(copyv), cap(copyv))  //[1 4 5 6] len=4,cap=6

//原始切片底层数组会用最后几位(删除的几位)补齐
fmt.Println(v2,len(v2), cap(v2)) //[1 2 5 6 5 6] len=6 cap=6

//删除后的切片不是新切片,修改会响应源数组
copyv[0] = 100
fmt.Println(copyv)  //[100 4 5 6]
fmt.Println(v2)  //[100 4 5 6 5 6]
}

func main() {
	s := make([]byte, 0 ,1000)
	var src = [10]byte {0,1,2,3,4,5,6,7,8,9}
	s = append(s , src[:]...)
	fmt.Println(len(src),cap(src),src,s)
	fmt.Printf("address of s %p\n",s)
	fmt.Println("llllllllllllllllllllllllllllll")
	start := 3
	size := 4
	s = utils.DeleteElementsFromSlice(s,start,size)
	fmt.Println(len(s),cap(s),s)
	fmt.Printf("address of s %p\n",s)
}