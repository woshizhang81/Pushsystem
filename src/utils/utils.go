package utils

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"os"
	"math/rand"
	"time"
	"strconv"
)

/*
获得配置文件路径
*/
func GetConfigPath() string {
	dir,_ := os.Getwd()
	config := dir + "/conf/"
	return config
}

/*
获得log文件夹路径
*/
func GetLogPath() string {
	dir,_ := os.Getwd()
	logPath := dir + "/log/"
	return logPath
}

/*
	设备唯一ID
*/
func UniqueId(idType int32 ,deviceId string) string {
	str := string(idType) + deviceId
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

/*
	得到字符串的哈希值
*/
func HasCode(s string ) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	return 0
}


/*
	delete 切片 中的固定长度的字节数组
*/
func DeleteElementsFromSlice(src []byte ,start int ,size int) []byte {
	sliceLen := len(src)
	if  (start + size) > sliceLen {
		fmt.Println("切片删除错误")
		os.Exit(1)
	}else{
	//	fmt.Println(start,size,"之间被丢掉的帧",src[start:start+size])
		src = append(src[:start], src[(start + size):]...)
	//	fmt.Println(start,size,"之间被丢掉之后的帧",src[:])
		return src
	}
	return src
}

/*
	返回 找到的 最近的子数组下标
*/
func FindSubByteArray(src []byte , needle []byte) int {
	var i int
	var j int
	srcLength := len(src)
	needleLength := len(needle)
	flag := true
	for i = 0 ; i < srcLength ; i++  {
		flag = true
		if (i + needleLength) > srcLength {
			return -1
		}
		for j = 0; j < needleLength ;j++  {
			if src[i + j] != needle[j]   {
				flag = false
			}
		}
		if flag {
			return i
		}
	}
	return -1
}

//生成6位随机数
func CreateCaptcha() int {
	str := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	num,err := strconv.Atoi(str)
	if err == nil {
		return num
	}else{
		return -1
	}
}

/*删除切片中固定的元素*/
func DeleteValueFormSlice(slice []string, key string ) bool {
	length := len(slice)
	for i := 0; i < length ; i++ {
		if slice[i] == key {
			slice = append(slice[:i], slice[i+1:]...)
			return true
		}
	}
	return false
}
