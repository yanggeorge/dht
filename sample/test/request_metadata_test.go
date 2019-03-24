package test

import (
	"fmt"
	"log"
	"net"
	"strings"
	"testing"
)

func TestSend(t *testing.T) {
	dial, err := net.Dial("tcp", ":8081")
	if err != nil {
		// handle error
	}
	conn := dial.(*net.TCPConn)
	i, err := sendData(conn, []byte("ym"))
	if err != nil || i != len("ym") {
		fmt.Println(err)
	}
	data, err := recvData(conn, 5)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}

func TestSever(t *testing.T) {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		// handle error
	}
	for {
		dial, err := ln.Accept()
		fmt.Println("connected")
		if err != nil {
			// handle error
		}
		conn := dial.(*net.TCPConn)
		go func(conn *net.TCPConn) {
			data, err := recvData(conn, 1024)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(data))
			conn.Write([]byte("hello"))

		}(conn)
	}
}

func TestRequestMetadata(t *testing.T) {
	infoHashString := "e84213a794f3ccd890382a54a64ca68b7e925433"
	ip := GetOutboundIP()
	ip4 := ip.To4().String()
	log.Printf("current ip =%s \n", ip4)
	addr := strings.Join([]string{ip4, "40959"}, ":")

	metadata, err := requestMetaData(infoHashString, addr)
	if err != nil {
		fmt.Println(err)
	}
	b := NewBDecode(metadata)
	dic, err := b.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(dic["name"])
}
