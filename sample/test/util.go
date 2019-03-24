package test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

func sendData(conn *net.TCPConn, data []byte) (int, error) {
	conn.SetWriteDeadline(time.Now().Add(time.Second * 50))
	return conn.Write(data)
}

func recvData(conn *net.TCPConn, n int) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 50))
	buf := new(bytes.Buffer)
	i, err := io.CopyN(buf, conn, int64(n))
	if err != nil {
		return nil, err
	}
	if int(i) != n {
		return nil, errors.New("cannot get size " + strconv.Itoa(n))
	}
	return buf.Bytes(), nil
}

func encodeToString(dat [20]byte) string {
	return hex.EncodeToString(dat[0:20])
}

func decodeString(s string) []byte {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return bytes
}

func getPeerId() []byte {
	data := make([]byte, 20)
	_, err := rand.Read(data)
	if err != nil {
		panic("error generate")
	}
	return data
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
