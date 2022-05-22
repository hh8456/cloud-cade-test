package base_net

import (
	"bytes"
	"cloud-cade-test/base-library/iface"
	"encoding/binary"
	"errors"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return 4
}

func (dp *DataPack) Pack(msg iface.IMessage) ([]byte, error) {
	dataBuf := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuf, binary.BigEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuf, binary.BigEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

func (dp *DataPack) Unpack(binaryData []byte, maxDataLen uint32) (iface.IMessage, error) {
	dataBuf := bytes.NewReader(binaryData)

	msg := &Message{}

	if err := binary.Read(dataBuf, binary.BigEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if msg.DataLen <= maxDataLen {
		if uint32(len(binaryData)) == dp.GetHeadLen()+msg.DataLen {
			msg.Data = binaryData[dp.GetHeadLen():]
			return msg, nil
		} else {
			return nil, errors.New("unpack fail, msgLen error")
		}
	} else {
		/*if maxDataLen < msg.DataLen */
		return nil, errors.New("too large msg data recieved")
	}
}

func (dp *DataPack) SetMsgLen(binaryData []byte, msgLen uint32) {
	if uint32(len(binaryData)) >= dp.GetHeadLen() {
		binary.BigEndian.PutUint32(binaryData, msgLen)
	}
}

func (dp *DataPack) UnpackMsgLen(binaryData []byte) (uint32, bool) {
	if uint32(len(binaryData)) < dp.GetHeadLen() {
		return 0, false
	}

	return binary.BigEndian.Uint32(binaryData), true
}
