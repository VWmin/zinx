package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/vwmin/zinx/utils"
	"github.com/vwmin/zinx/ziface"
)

/**
感觉像是一个工具类而已
*/
type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包的长度的方法
func (dataPack *DataPack) GetHeadLen() uint32 {
	// DataLen(uint32 4bytes) + Id(uint32 4bytes)
	return 4 + 4
}

// 封包方法
func (dataPack *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	// 将DataId写入buf
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 将DataLen写入buf
	if err := binary.Write(buf, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 将Data写入Buf
	if err := binary.Write(buf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 拆包方法，获得消息头部信息
func (dataPack *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	reader := bytes.NewReader(binaryData)

	var msgId uint32
	var dataLen uint32

	/* 读消息id */
	if err := binary.Read(reader, binary.LittleEndian, &msgId); err != nil {
		return nil, err
	}

	/* 读消息长度 */
	if err := binary.Read(reader, binary.LittleEndian, &dataLen); err != nil {
		return nil, err
	}

	limit := utils.GlobalObject.MaxPackageSize

	// dataLen超出长度
	if limit > 0 && dataLen > limit {
		return nil, errors.New("package oversize error")
	}

	// 如果不传入conn，需要由调用者判断数据长度，再从conn读取相应内容
	//if err := binary.Read(reader, binary.LittleEndian, data); err != nil {
	//	return nil, err
	//}

	return NewMessageHead(msgId,dataLen), nil
}
