package dht

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"testing"
	"time"
)

func TestFetchMetaData(t *testing.T) {
	var (
		length       int
		msgType      byte
		piecesNum    int
		pieces       [][]byte
		utMetadata   int
		metadataSize int
	)
	infoHash, err := hex.DecodeString("e84213a794f3ccd890382a54a64ca68b7e925433")
	if err != nil {
		return
	}
	address := "192.168.0.103:40959"

	dial, err := net.DialTimeout("tcp", address, time.Second*50)

	conn := dial.(*net.TCPConn)
	conn.SetLinger(0)
	defer conn.Close()

	data := bytes.NewBuffer(nil)
	data.Grow(BLOCK)
	//n, err := io.CopyN(data, conn, int64(68))
	//if err != nil || n != int64(68) {
	//	fmt.Println(err)
	//	return
	//}
	if sendHandshake(conn, infoHash, []byte(randomString(20))) != nil ||
		read(conn, 68, data) != nil ||
		onHandshake(data.Next(68)) != nil ||
		sendExtHandshake(conn) != nil {
		return
	}

	for {
		length, err = readMessage(conn, data)
		if err != nil {
			return
		}

		if length == 0 {
			continue
		}

		msgType, err = data.ReadByte()
		if err != nil {
			return
		}
		fmt.Printf("msgType=%d, length=%d \n", msgType, length)
		switch msgType {
		case EXTENDED:
			extendedID, err := data.ReadByte()
			if err != nil {
				return
			}

			payload, err := ioutil.ReadAll(data)
			if err != nil {
				return
			}

			if extendedID == 0 {
				if pieces != nil {
					return
				}

				utMetadata, metadataSize, err = getUTMetaSize(payload)
				fmt.Println(utMetadata)
				if err != nil {
					return
				}

				piecesNum = metadataSize / BLOCK
				if metadataSize%BLOCK != 0 {
					piecesNum++
				}

				pieces = make([][]byte, piecesNum)
				//go wire.requestPieces(conn, utMetadata, metadataSize, piecesNum)

				continue
			}

			if pieces == nil {
				return
			}

			d, index, err := DecodeDict(payload, 0)
			if err != nil {
				return
			}
			dict := d.(map[string]interface{})

			if err = ParseKeys(dict, [][]string{
				{"msg_type", "int"},
				{"piece", "int"}}); err != nil {
				return
			}

			if dict["msg_type"].(int) != DATA {
				continue
			}

			piece := dict["piece"].(int)
			pieceLen := length - 2 - index

			if (piece != piecesNum-1 && pieceLen != BLOCK) ||
				(piece == piecesNum-1 && pieceLen != metadataSize%BLOCK) {
				return
			}

			pieces[piece] = payload[index:]

			//if wire.isDone(pieces) {
			//	metadataInfo := bytes.Join(pieces, nil)
			//
			//	info := sha1.Sum(metadataInfo)
			//	if !bytes.Equal(infoHash, info[:]) {
			//		return
			//	}
			//
			//	wire.responses <- Response{
			//		Request:      r,
			//		MetadataInfo: metadataInfo,
			//	}
			//	return
			//}
		default:
			data.Reset()
		}
	}
}
