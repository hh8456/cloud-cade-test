package main

import (
	"bufio"
	"cloud-cade-test/base-library/base_net"
	"flag"
	"fmt"
	"os"
)

func recvLoop(s *base_net.Socket) {
	defer println("exit recvLoop")
	for {
		msg, e := s.ReadOne()
		if e != nil {
			return
		}
		println(string(msg[4:]))
	}
}

func sendLoop(chanSend chan []byte, s *base_net.Socket) {
	defer println("exit sendLoop")

	for {
		msg, ok := <-chanSend
		if ok {
			e := s.Send(msg)
			if e != nil {
				return
			}
		} else {
			return
		}
	}
}

func main() {
	addr := flag.String("a", "127.0.0.1:4567", "server ip")
	flag.Parse()
	s, e := base_net.ConnectSocket(*addr, 4096)
	if e != nil {
		fmt.Printf("connect fail, error: %v, please reset client", e)
		return
	}

	chanSend := make(chan []byte, 20)

	go sendLoop(chanSend, s)
	go recvLoop(s)

	println("input login name( must be english ):")
	for {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadBytes('\n')

		//println(string(line))
		if err != nil {
			fmt.Println(err.Error())
		}

		dp := base_net.NewDataPack()
		buf, _ := dp.Pack(base_net.NewMsgPackage(line[:len(line)-1])) // 去掉 '\n'
		chanSend <- buf
	}

}
