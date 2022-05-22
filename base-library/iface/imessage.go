package iface

type IMessage interface {
	GetDataLen() uint32
	GetData() []byte

	SetData([]byte)
}
