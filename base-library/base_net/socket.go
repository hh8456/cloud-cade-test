package base_net

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	defTimeoutDuration = time.Minute
)

var ErrOverMaxReadingSize = errors.New("over max reading size")

type Socket struct {
	conn                 net.Conn
	recvBuf              []byte
	reader               *bufio.Reader
	maxReadSize          uint32
	timeoutReadDuration  time.Duration
	timeoutWriteDuration time.Duration
}

// 包体最大长度为 maxReadSize, 不算包头
func CreateSocket(conn net.Conn, maxReadSize uint32) *Socket {
	dp := NewDataPack()
	return &Socket{
		conn:                 conn,
		recvBuf:              make([]byte, dp.GetHeadLen()),
		reader:               bufio.NewReader(conn),
		maxReadSize:          maxReadSize,
		timeoutReadDuration:  60 * defTimeoutDuration,
		timeoutWriteDuration: defTimeoutDuration,
	}
}

func (s *Socket) Conn() net.Conn {
	return s.conn
}

func (s *Socket) ReadOne() ([]byte, error) {
	b, e := s.read()
	if e != nil {
		return nil, e
	}
	return b, nil
}

func (s *Socket) read() ([]byte, error) {
	s.conn.SetReadDeadline(time.Now().Add(s.timeoutReadDuration))
	if _, err := io.ReadFull(s.reader, s.recvBuf); err != nil {
		return nil, err
	}

	dp := NewDataPack()
	msgLen, _ := dp.UnpackMsgLen(s.recvBuf)

	if msgLen > s.maxReadSize {
		return nil, fmt.Errorf("read too large( len = %d ) data ",
			msgLen)
	}

	s.conn.SetReadDeadline(time.Now().Add(s.timeoutReadDuration))
	length := dp.GetHeadLen() + msgLen
	buf := make([]byte, length)

	copy(buf, s.recvBuf)
	if _, err := io.ReadFull(s.reader, buf[dp.GetHeadLen():length]); err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *Socket) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *Socket) Close() {
	s.conn.Close()
}

func (s *Socket) Send(msg []byte) error {
	_, e := s.conn.Write(msg)
	return e
}

func (s *Socket) SendPbMsg(pbMsg proto.Message) error {
	if pbMsg != nil {
		msg, err := proto.Marshal(pbMsg)
		if err != nil {
			return err
		}

		return s.SendPbBuf(msg)

	} else {
		return s.SendPbBuf(nil)
	}
}

func (s *Socket) SendPbBuf(msg []byte) error {
	dp := NewDataPack()
	buf, err := dp.Pack(NewMsgPackage(msg))
	if err != nil {
		return err
	}

	_, e := s.conn.Write(buf)
	return e
}
