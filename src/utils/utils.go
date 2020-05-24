package utils

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"os"
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
		src = append(src[:start], src[(start + size):]...)
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
		if (i + needleLength) >= srcLength {
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

