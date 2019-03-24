package test

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	BLOCK = 16384
)

func packHandshake(infoHash string) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int8(19))
	if err != nil {
		panic(err)
	}
	head := []byte("BitTorrent protocol")
	err = binary.Write(buf, binary.BigEndian, head)
	if err != nil {
		panic(err)
	}
	reserved := []byte("\x00\x00\x00\x00\x00\x10\x00\x01")

	err = binary.Write(buf, binary.BigEndian, reserved)
	if err != nil {
		panic(err)
	}

	err = binary.Write(buf, binary.BigEndian, decodeString(infoHash))
	if err != nil {
		panic(err)
	}
	err = binary.Write(buf, binary.BigEndian, getPeerId())
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func unpackHandshake(data []byte) (reserved []byte, infoHash []byte, peerId []byte, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, int8(19))
	if err != nil {
		return nil, nil, nil, err
	}
	head := []byte("BitTorrent protocol")
	err = binary.Write(buf, binary.BigEndian, head)
	if err != nil {
		return nil, nil, nil, err
	}
	prefix := buf.Bytes()
	if !bytes.HasPrefix(data, prefix) {
		return nil, nil, nil, fmt.Errorf("not begin")
	}
	data = data[20:]
	reserved = data[:8]
	infoHash = data[8:28]
	peerId = data[28:48]
	return reserved, infoHash, peerId, nil
}

func packExtend(dic map[string]interface{}, id, extId int) ([]byte, error) {
	data, err := encodeToBin(dic)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	msgLength := len(data)
	err = writeInt32(buf, msgLength+2)
	if err != nil {
		return nil, err
	}
	err = writeInt8(buf, id)
	if err != nil {
		return nil, err
	}
	err = writeInt8(buf, extId)
	if err != nil {
		return nil, err
	}
	err = writeBytes(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func requestMetaData(infoHashString, addr string) (metadata []byte, err error) {
	dial, err := net.DialTimeout("tcp", addr, time.Second*15)
	if err != nil {
		return
	}
	var (
		pieces      [][]byte
		pieceNum    int
		ut_metadata int
	)

	conn := dial.(*net.TCPConn)
	conn.SetLinger(0)
	defer conn.Close()

	data := packHandshake(infoHashString)
	i, err := sendData(conn, data)
	if err != nil {
		return
	}

	if i != 68 {
		return nil, errors.New("not 68 bytes")
	}
	_, err = recvData(conn, 68)
	if err != nil {
		return
	}

	data, err = getExtendPack()
	if err != nil {
		return nil, err
	}
	_, err = sendData(conn, data)
	if err != nil {
		return nil, err
	}

	for {
		tmp, err := recvData(conn, 4)

		if err != nil {
			return nil, err
		}

		if len(tmp) == 0 {
			continue
		}

		msg_len := binary.BigEndian.Uint32(tmp)
		tmp, err = recvData(conn, 1)
		if err != nil {
			return nil, err
		}
		id := int8(tmp[0])

		var metadataSize int
		switch id {
		case 20:
			tmp, err = recvData(conn, 1)
			if err != nil {
				return nil, err
			}
			ext_id := int(tmp[0])
			data, err = recvData(conn, int(msg_len)-2)
			if err != nil {
				return nil, err
			}
			switch ext_id {
			case 0:
				// extend handshake
				b := NewBDecode(data)
				dic, err := b.Dic()
				if err != nil {
					return nil, err
				}

				metadataSize = dic["metadata_size"].(int)

				tail := metadataSize % BLOCK
				if tail == 0 {
					pieceNum = metadataSize / BLOCK
				} else {
					pieceNum = int(metadataSize/BLOCK) + 1
				}

				pieces = make([][]byte, pieceNum)
				ut_metadata = dic["m"].(map[string]interface{})["ut_metadata"].(int)
				go requestPieces(conn, pieceNum, ut_metadata)

			case 1:
				// ut_metadata
				b := NewBDecode(data)
				dic, err := b.Dic()
				if err != nil {
					return nil, err
				}
				pieces[dic["piece"].(int)] = data[b.i:]

				if isDone(pieces) {
					metadataInfo := bytes.Join(pieces, nil)
					infoHash := sha1.Sum(metadataInfo)
					if hex.EncodeToString(infoHash[0:20]) != infoHashString {
						return nil, errors.New("infoHash not equal")
					}
					return metadataInfo, nil
				}

			default:

			}
		case 9:
			// port
			tmp, err = recvData(conn, int(msg_len)-1)
			if err != nil {
				return nil, err
			}
		case 1:
		default:
			_, err := recvData(conn, int(msg_len)-1)
			if err != nil {
				return nil, err
			}
		}
	}
}

func isDone(pieces [][]byte) bool {
	for i := 0; i < len(pieces); i++ {
		if pieces[i] == nil {
			return false
		}
	}
	return true
}

func requestPieces(conn *net.TCPConn, pieceNum int, ut_metadata int) {
	for i := 0; i < pieceNum; i++ {
		requestBytes, err := getPieceRequestPack(i, ut_metadata)
		if err != nil {
			return
		}
		_, err = sendData(conn, requestBytes)
		if err != nil {
			return
		}
	}
}

func getExtendPack() ([]byte, error) {
	dic := make(map[string]interface{})
	metadic := make(map[string]interface{})
	metadic["ut_metadata"] = 1
	dic["m"] = metadic

	return packExtend(dic, 20, 0)
}

func getPieceRequestPack(piece, ext_id int) ([]byte, error) {
	dic := make(map[string]interface{})
	dic["msg_type"] = 0
	dic["piece"] = piece
	return packExtend(dic, 20, ext_id)

}
