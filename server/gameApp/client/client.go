package client

import (
	"cloud-cade-test/base-library/base_net"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	serialNumber int64
)

type Client struct {
	*base_net.Connection
	lock        sync.RWMutex
	Name        string
	loginClient func(name string) bool
	//mapLogicFunc map[string]func(*Client, ...string)
	mapLogicFunc map[string]func(string, ...string) []byte
}

func NewClient(c *base_net.Socket,
	loginClient func(string) bool,
	mapLogicFunc map[string]func(string, ...string) []byte,
) *Client {
	cli := &Client{}

	conn := base_net.NewConnection(c, atomic.AddInt64(&serialNumber, 1), 20, 20, cli.recvData, cli.exit)
	cli.Connection = conn
	cli.loginClient = loginClient
	cli.mapLogicFunc = map[string]func(string, ...string) []byte{}
	// 虽然并发读取 map 可以不用加锁,出于健壮性考虑,还是让每个 client 有自己的 map 对象
	for k, v := range mapLogicFunc {
		cli.mapLogicFunc[k] = v
	}

	return cli
}

func (c *Client) Run() {
	c.Connection.Run()
}

func (c *Client) exit() {
	fmt.Printf("client disconnect, serialNumber: %d\n", c.ConnId())
}

func (c *Client) recvData(binData []byte) {
	dp := base_net.NewDataPack()
	msg, _ := dp.Unpack(binData, 4096) // 客户端发来的包,不能超过 4k
	fmt.Printf("recv client msg, msgLen: %d, msgContext: %s\n", msg.GetDataLen(), msg.GetData())

	fields := strings.Fields(string(msg.GetData()))
	logicFunc, ok := c.mapLogicFunc[fields[0]]
	if ok {
		// 处理 /createAlliance /allianceList /joinAlliance /dismissAlliance
		// /increaseCapacity /storeItem /destroyItem /clearup
		if len(fields) == 1 {
			msg := logicFunc(c.Name, "")
			buf, _ := dp.Pack(base_net.NewMsgPackage(msg))
			c.Send(buf)
		} else {
			msg := logicFunc(c.Name, fields[1:]...)
			buf, _ := dp.Pack(base_net.NewMsgPackage(msg))
			println(string(msg))
			c.Send(buf)
		}

	} else {
		if c.Name != "" {
			buf, _ := dp.Pack(base_net.NewMsgPackage([]byte("unrecognized command")))
			c.Send(buf)
			return
		}

		name := string(binData)
		if c.loginClient(name) {
			// 登录成功
			c.Name = name
			buf, _ := dp.Pack(base_net.NewMsgPackage([]byte("login success")))
			c.Send(buf)
		} else {
			// 名字重复, 需要重试
			buf, _ := dp.Pack(base_net.NewMsgPackage([]byte("the name is already in use, enter a new name")))
			c.Send(buf)
		}
	}
}
