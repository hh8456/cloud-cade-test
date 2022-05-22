package base_net

import (
	"cloud-cade-test/base-library/iface"
	"context"
	"encoding/binary"
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/proto"
)

type Connection struct {
	*Socket
	connId      int64 // 连接id, 全局唯一
	lock        sync.RWMutex
	isClosed    bool
	isRunning   uint32
	ctx         context.Context
	cancel      context.CancelFunc
	chanSendMsg chan []byte // 保存将要发给 socket 的数据
	chanRecvMsg chan []byte // 保存从 socket 读到的数据
	handleFunc  func([]byte)
	done        func()
	mapProperty map[string]interface{}
}

func NewConnection(s *Socket, connId int64,
	sendChannelSize, recvChannelSize uint32,
	handleFunc func([]byte), done func()) *Connection {
	ctx, cancel := context.WithCancel(context.Background())
	return &Connection{Socket: s, connId: connId,
		ctx: ctx, cancel: cancel,
		chanSendMsg: make(chan []byte, sendChannelSize),
		chanRecvMsg: make(chan []byte, recvChannelSize),
		handleFunc:  handleFunc,
		done:        done,
		mapProperty: map[string]interface{}{}}
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.mapProperty[key] = value
}

func (c *Connection) GetProperty(key string) interface{} {
	if v, find := c.mapProperty[key]; find {
		return v
	}

	return nil
}

func (c *Connection) ConnId() int64 {
	return c.connId
}

func (c *Connection) Run() {
	if 1 == atomic.AddUint32(&c.isRunning, 1) {
		go c.recvLoop()
		go c.sendLoop()
		go c.handleLoop()
	}
}

// 用于通知外部协程关闭
func (c *Connection) Done() <-chan struct{} {
	return c.ctx.Done()
}

// recvLoop 和 sendLoop 其中一方退出时, 会触发另外一方退出
func (c *Connection) recvLoop() {

	for {
		binData, err := c.ReadOne()
		if err != nil {
			c.Close()
			return
		}

		select {
		case c.chanRecvMsg <- binData:

		default:
		}
	}
}

// recvLoop 和 sendLoop 其中一方退出时, 会触发另外一方退出
func (c *Connection) sendLoop() {

	for {
		select {
		case binData, ok := <-c.chanSendMsg:
			if ok {
				c.Socket.Send(binData)
			}

		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Connection) handleLoop() {
	for {
		select {
		case binMsg, ok := <-c.chanRecvMsg:
			if ok {
				c.handleFunc(binMsg)
			}

		case <-c.ctx.Done():
			c.done()
			return
		}
	}
}

func (c *Connection) SendPb(pb proto.Message) error {
	buf, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	dp := DataPack{}
	binData := make([]byte, dp.GetHeadLen()+uint32(len(buf)))
	binary.BigEndian.PutUint32(binData, uint32(len(buf)))
	copy(binData[dp.GetHeadLen():], buf)
	c.Send(binData)
	return nil
}

// Connection.Send 是把 binData 投递到发送缓冲区
// 如果需要直接发送, 就使用 Connection.Socket.Send
func (c *Connection) Send(binData []byte) {
	// 加锁是为了防止写 closed c.chanSendMsg 崩溃
	// go test bench 的并发压测显示, 不加锁 40 ns/op, 加读锁后 62 ns/op
	c.lock.RLock()
	if c.isClosed == false {
		select {
		case c.chanSendMsg <- binData:

		default:
		}

	}
	c.lock.RUnlock()
}

func (c *Connection) SendMessage(message iface.IMessage) {
	dp := NewDataPack()
	buf, err := dp.Pack(message)
	if err != nil {
		return
	}

	c.Send(buf)
}

func (c *Connection) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.isClosed == true {
		return
	}

	c.isClosed = true
	c.Socket.Close()
	c.cancel()
}
