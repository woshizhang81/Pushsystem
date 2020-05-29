package protocol

import (
	"Pushsystem/src/utils"
	"encoding/binary"
	//"fmt"
	"container/list"
	//"os"
	//"os"
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

/*标准化处理 0xff 0xff 开头*/
func (obj *ProtoCheck)StandardDealMent()  bool {
	fixHead := [2]byte{0xff,0xff}
	index := utils.FindSubByteArray(obj.storeBuf, fixHead[:])
	if index == -1 { //说明没找到头 ，该帧需要丢掉，一般不可能出现这种情况
		if obj.storeBuf[len(obj.storeBuf) -1] == 0xff {
			//fmt.Println("无0xff丢掉的帧",obj.storeBuf[:len(obj.storeBuf) -1])
			obj.storeBuf = append(obj.storeBuf[0:0], 0xff)
			return false
		}else {
			//fmt.Println("无0xff丢掉的帧",obj.storeBuf)
			obj.storeBuf = obj.storeBuf[0:0] //清空该buf
			return false
		}
	} else {
		if index != 0 {
			//丢掉index之前的所有帧
			//fmt.Println("index之前的所有帧",obj.storeBuf[:index])
			retBuf := utils.DeleteElementsFromSlice(obj.storeBuf,0,index)
			obj.storeBuf = retBuf
		}
	}
	return true
}
/*
	list 保存了解出来的帧
*/
func (obj *ProtoCheck)CheckAndGetProtocolBuffer(bufRcv []byte,list *list.List ) bool  {
	//fmt.Println("before connect", obj.storeBuf)
	obj.storeBuf = append(obj.storeBuf, bufRcv[:]...)
	ret := obj.StandardDealMent()
	if ret {
		obj.decodeFromBinaryStream(list)
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
	if curBufLen < 4 { //说明当前buf 不包含 包的字节数
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
				//fmt.Println("校验成功的帧",len(obj.storeBuf[:packageSize]),obj.storeBuf[:packageSize])
				//异步操作，无法保证 同一会影响切片的 数据
				// 复制一个切片 //复制会导致效率低，暂时保留找别的解决方法？？
				var frame []byte
				frame = append(frame[0:0],obj.storeBuf[:packageSize]...)
				list.PushBack(frame[:])
			}else {
				//fmt.Println("校验失败的帧",packageSize ,obj.storeBuf[:packageSize])
			}
			//丢掉一帧 storeBuf 里的数据
			//obj.storeBuf = append(obj.storeBuf[0:0], obj.storeBuf[packageSize:]...)
			obj.storeBuf = utils.DeleteElementsFromSlice(obj.storeBuf,0,int(packageSize))
		//	fmt.Println("cur buf frame",len(obj.storeBuf),obj.storeBuf)
			if curBufLen > packageSize {
				ret := obj.StandardDealMent()
				if ret {
					obj.decodeFromBinaryStream(list) //递归一次继续找
				}
			}
			return true //都认为返回成功了
		}

	}
}

