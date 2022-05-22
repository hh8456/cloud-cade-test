package base_net

import (
	"fmt"
	"net"
	"time"
)

type Listener struct {
	isClosed bool
	listener net.Listener
}

func CreateListener() *Listener {
	return &Listener{}
}

func (l *Listener) Start(addr string, onEstablish func(c net.Conn)) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	l.listener = ln
	var tempDelay time.Duration // how long to sleep on accept failure

	for {
		conn, err := ln.Accept()
		if l.isClosed {
			break
		}
		if err != nil {
			str := fmt.Sprintf("监听组件出现错误: %v", err)
			println(str)
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0
		onEstablish(conn)
	}

	return nil
}

func (l *Listener) Close() error {
	l.isClosed = true
	return l.listener.Close()
}

func ConnectSocket(addr string, maxReadSize uint32) (*Socket, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return CreateSocket(conn, maxReadSize), nil
}
