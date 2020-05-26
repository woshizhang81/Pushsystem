package protocol

import (
	"Pushsystem/src/utils"
	"encoding/binary"
	"fmt"
)

/*
	检查并从buf中解出 标准协议帧 考虑粘包的问题
	输入 切片为拼接后的 切片
*/

const DefaultCapLength = 4096
//得到完整的一帧数据
type CallBackOnGetStream func (handle *ProtoCheck , stream []byte)


type ProtoCheck struct{
	Handle   interface{} //句柄 相关对象的
	storeBuf []byte	   //负责包的拼接
	//CallBack CallBackOnGetStream
}

func (obj *ProtoCheck)Init()  {
	obj.storeBuf = make([]byte ,DefaultCapLength, DefaultCapLength)
}

func (obj *ProtoCheck)SetCallBack(cbfun CallBackOnGetStream,handle interface{}) {
	//obj.CallBack = cbfun
	obj.Handle = handle
}


/*
	list 保存了解出来的帧
*/
func (obj *ProtoCheck)CheckAndGetProtocolBuffer(bufRcv []byte,allbuffer [][]byte) bool  {
	//找到 头为0xff 0xff 的起始地址
	obj.storeBuf = append(obj.storeBuf, bufRcv[:]...)
	fixHead := [2]byte{0xff,0xff}
	index := utils.FindSubByteArray(obj.storeBuf, fixHead[:])
	if index == -1 { //说明没找到头 ，该帧需要丢掉，一般不可能出现这种情况
		if obj.storeBuf[len(obj.storeBuf) -1] == 0xff {
			obj.storeBuf = append(obj.storeBuf[0:0], 0xff)
			return false
		}else {
			fmt.Println(obj.storeBuf)
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
			obj.decodeFromBinaryStream(allbuffer)

		}
	}

	return true
}

/*
	从数据流中解码按规则解码
	headSize  offset (2)
	PackType  offset (2 + headSize)
	PackSize  offset (2 + headSize + 1)
*/
const HeadSizeOffset = 2
const PackageSizeBytes  = 2

func (obj *ProtoCheck) decodeFromBinaryStream(allbuffer [][]byte)  bool {
	//func (obj *ProtoCheck) decodeFromBinaryStream(list *list.List)  bool {
	curBufLen := uint16(len(obj.storeBuf))
	if curBufLen < 3 { //说明当前buf 不包含 Head中的 HeadSize
		return false
	}else{
		//此时obj.storeBuf 开头必为0xff 0xff
		headSize := uint16(obj.storeBuf[HeadSizeOffset])
		packSizeOffset := HeadSizeOffset + headSize + 1
		if curBufLen < packSizeOffset + 1 + PackageSizeBytes { //说明长度不包含包总字节数目
			return false
		}else {
			//得到 包总长度
			totalPackageSize := binary.BigEndian.Uint16(obj.storeBuf[packSizeOffset:2])
			if curBufLen < totalPackageSize {
				//说明当前长度小于 包总长度
				return false
			} else {
				crc16 := binary.BigEndian.Uint16(obj.storeBuf[(totalPackageSize - 3) : (totalPackageSize - 1)])
				crcCalculate := utils.Crc16(obj.storeBuf[:(totalPackageSize -3)])
				if crc16 == crcCalculate {
					//校验成功了
					allbuffer = append(allbuffer, obj.storeBuf[:totalPackageSize -1])
					//list.PushBack(obj.storeBuf[:totalPackageSize -1 ])
				/*	if obj.CallBack != nil {
						obj.CallBack(obj,obj.storeBuf[:totalPackageSize -1])
					}else{
						fmt.Println("未设置回调函数")
						os.Exit(1)  //程序异常
					}
				 */
				}else {
					fmt.Println("校验失败")
				}
				if curBufLen > totalPackageSize {
					obj.storeBuf = utils.DeleteElementsFromSlice(obj.storeBuf,0,int(totalPackageSize))
					ret := obj.decodeFromBinaryStream(allbuffer) //递归一次继续找
					if !ret { //结束条件 失败
						return false
					}
				}
				return true //都认为返回成功了
			}
		}
	}
}

