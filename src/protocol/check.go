package protocol

import (
	"Pushsystem/src/utils"
	"encoding/binary"
	"fmt"
	"container/list"
)

/*
	检查并从buf中解出 标准协议帧 考虑粘包的问题
	输入 切片为拼接后的 切片
*/

const DefaultCapLength = 4096
//得到完整的一帧数据

//type OnGetFrame func (handle * interface{},buf []byte)

type ProtoCheck struct{
	storeBuf []byte	   //负责包的拼接
//	handle interface{}
//	callback OnGetFrame
}

func (obj *ProtoCheck)Init()  {
	obj.storeBuf = make([]byte , 0 , DefaultCapLength)
}
/*
	list 保存了解出来的帧
*/
//func (obj *ProtoCheck)CheckAndGetProtocolBuffer(bufRcv []byte,allbuffer [][]byte) bool  {
	func (obj *ProtoCheck)CheckAndGetProtocolBuffer(bufRcv []byte,list *list.List ) bool  {
	//找到 头为0xff 0xff 的起始地址
	obj.storeBuf = append(obj.storeBuf, bufRcv[:]...)
	fixHead := [2]byte{0xff,0xff}
	index := utils.FindSubByteArray(obj.storeBuf, fixHead[:])
	if index == -1 { //说明没找到头 ，该帧需要丢掉，一般不可能出现这种情况
		if obj.storeBuf[len(obj.storeBuf) -1] == 0xff {
			obj.storeBuf = append(obj.storeBuf[0:0], 0xff)
			return false
		}else {
			//fmt.Println(obj.storeBuf)
			obj.storeBuf = obj.storeBuf[0:0] //清空该buf
			return false
		}
	} else {
		if index != 0 {
			//丢掉index之前的所有帧
			retBuf := utils.DeleteElementsFromSlice(obj.storeBuf,0,index)
			obj.storeBuf = retBuf
		}else {
			//一般情况在这里
			//开始按规则截取包
			obj.decodeFromBinaryStream(list)
		}
	}

	return true
}

/*
	从数据流中解码按规则解码
	PackSize  offset (2 + headSize + 1)
*/
const packageSizeOffset = 2

//func (obj *ProtoCheck) decodeFromBinaryStream(allbuffer [][]byte)  bool {
	func (obj *ProtoCheck) decodeFromBinaryStream(list *list.List)  bool {
	curBufLen := uint16(len(obj.storeBuf))
	if curBufLen < 2 { //说明当前buf 不包含 包的字节数
		return false
	}else{
		//此时obj.storeBuf 开头必为0xff 0xff
		packageSize := binary.BigEndian.Uint16(obj.storeBuf[packageSizeOffset:])
		if packageSize > curBufLen { //当前长度小于包长度
			return false
		} else {
			crc16 := binary.BigEndian.Uint16(obj.storeBuf[(packageSize - 2) : ])
			crcCalculate := utils.Crc16(obj.storeBuf[:(packageSize - 2)])

			if crc16 == crcCalculate {
			//校验成功了
			//	fmt.Println("curdatas",len(obj.storeBuf[:packageSize]),obj.storeBuf[:packageSize])
				//list = append(list, obj.storeBuf[:packageSize])
				list.PushBack(obj.storeBuf[:packageSize])
			}else {
				fmt.Println("校验失败")
			}

			if curBufLen > packageSize {
				obj.storeBuf = utils.DeleteElementsFromSlice(obj.storeBuf,0,int(packageSize))
				ret := obj.decodeFromBinaryStream(list) //递归一次继续找
				if !ret { //结束条件 list
					return false
				}
			}
			return true //都认为返回成功了
		}

	}
}

