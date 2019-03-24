package test

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"
	"testing"
	"time"
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

func TestPack(t *testing.T) {
	ip := GetOutboundIP()
	ip4 := ip.To4().String()
	log.Printf("current ip =%s \n", ip4)
	addr := strings.Join([]string{ip4, "40959"}, ":")
	dial, err := net.DialTimeout("tcp", addr, time.Second*50)
	check(err)

	var (
		pieces      [][]byte
		pieceNum    int
		ut_metadata int
	)

	conn := dial.(*net.TCPConn)
	conn.SetLinger(0)
	defer conn.Close()

	infoHashString := "e84213a794f3ccd890382a54a64ca68b7e925433"
	data := packHandshake(infoHashString)
	i, err := sendData(conn, data)
	check(err)
	if i != 68 {
		fmt.Printf("err")
	}
	resp, err := recvData(conn, 68)

	if err != nil {
		fmt.Println(err)
		return
	}
	reserved, infoHash, peerId, err := unpackHandshake(resp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("% x ,% x ,% x \n", reserved, infoHash, peerId)

	data, err = getExtendPack()
	if err != nil {
		fmt.Println(err)
	}
	_, err = sendData(conn, data)
	if err != nil {
		fmt.Println(err)
	}
	for {
		tmp, err := recvData(conn, 4)

		if err != nil {
			fmt.Println(err)
			return
		}

		if len(tmp) != 4 {
			continue
		}

		msg_len := binary.BigEndian.Uint32(tmp)
		tmp, err = recvData(conn, 1)
		if err != nil {
			fmt.Println(err)
		}
		id := int8(tmp[0])

		fmt.Printf("msg_len=%d, id=%d\n", msg_len, id)
		var metadata_size int
		switch id {
		case 20:
			tmp, err = recvData(conn, 1)
			if err != nil {
				fmt.Println(err)
				return
			}
			ext_id := int(tmp[0])
			data, err = recvData(conn, int(msg_len)-2)
			if err != nil {
				fmt.Println(err)
				return
			}
			switch ext_id {
			case 0:
				// extend handshake
				b := NewBDecode(data)
				dic, err := b.Dic()
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(dic)
				metadata_size = dic["metadata_size"].(int)
				fmt.Println(metadata_size)

				tail := metadata_size % BLOCK
				if tail == 0 {
					pieceNum = metadata_size / BLOCK
				} else {
					pieceNum = int(metadata_size/BLOCK) + 1
				}
				pieces = make([][]byte, pieceNum)
				ut_metadata = dic["m"].(map[string]interface{})["ut_metadata"].(int)
				go requestPieces(conn, pieceNum, ut_metadata)

			case 1:
				// ut_metadata
				b := NewBDecode(data)
				dic, err := b.Dic()
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(dic, b.i, b.n, b.n-b.i)
				pieces[dic["piece"].(int)] = data[b.i:]

				if isDone(pieces) {
					metadataInfo := bytes.Join(pieces, nil)
					b = NewBDecode(metadataInfo)

					infoHash := sha1.Sum(metadataInfo)
					if hex.EncodeToString(infoHash[0:20]) != infoHashString {
						return
					}
					dic, err := b.Dic()
					if err != nil {
						return
					}
					fmt.Println(dic["name"])
					return
				}

			default:

			}
		case 9:
			// port
			tmp, err = recvData(conn, int(msg_len)-1)
			if err != nil {
				fmt.Println(err)
				return
			}
			port := binary.BigEndian.Uint16(tmp)
			fmt.Printf("port=%d\n", port)
		case 1:
		default:
			recvData(conn, int(msg_len)-1)
		}
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
