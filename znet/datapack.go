package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	// head [id---len] id---4 len---4
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	databuff := bytes.NewBuffer([]byte{})
	// write datalen
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	// write msgid
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// write data
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return databuff.Bytes(), nil
}

func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	databuff := bytes.NewReader(binaryData)
	msg := &Message{}
	// read datalen
	if err := binary.Read(databuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// read msgid
	if err := binary.Read(databuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data received")
	}
	return msg, nil
}
