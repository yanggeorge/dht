package test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sort"
	"strconv"
)

func encodeToBin(dic map[string]interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := encodeDic(buf, dic)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encodeDic(buf *bytes.Buffer, dic map[string]interface{}) error {
	keys := make([]string, 0)
	err := writeByte(buf, 'd')
	if err != nil {
		return err
	}
	for k := range dic {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		err := encodeString(buf, k)
		if err != nil {
			return err
		}
		err = encodeElement(buf, dic[k])
		if err != nil {
			return err
		}

	}

	err = writeByte(buf, 'e')
	if err != nil {
		return err
	}
	return nil
}

func encodeElement(buf *bytes.Buffer, val interface{}) error {
	switch val.(type) {
	case string:
		s, ok := val.(string)
		if !ok {
			return errors.New("cannot cast to string")
		}
		err := encodeString(buf, s)
		if err != nil {
			return err
		}
	case int:
		i, ok := val.(int)
		if !ok {
			return errors.New("cannot cast to int")
		}
		err := encodeInt(buf, i)
		if err != nil {
			return err
		}
	case []interface{}:
		lis, ok := val.([]interface{})
		if !ok {
			return errors.New("cannot cast to []interface{}")
		}
		err := encodeList(buf, lis)
		if err != nil {
			return err
		}
	case map[string]interface{}:
		dic, ok := val.(map[string]interface{})
		if !ok {
			return errors.New("cannot cast to map[string]interface{}")
		}
		err := encodeDic(buf, dic)
		if err != nil {
			return err
		}
	default:

	}
	return nil
}

func encodeList(buf *bytes.Buffer, lis []interface{}) error {
	err := writeByte(buf, 'l')
	if err != nil {
		return err
	}

	for _, ele := range lis {
		err := encodeElement(buf, ele)
		if err != nil {
			return err
		}
	}

	err = writeByte(buf, 'e')
	if err != nil {
		return err
	}
	return nil
}

func encodeInt(buf *bytes.Buffer, i int) error {
	err := writeByte(buf, 'i')
	if err != nil {
		return err
	}
	err = writeBytes(buf, []byte(strconv.Itoa(i)))
	if err != nil {
		return err
	}
	err = writeByte(buf, 'e')
	if err != nil {
		return err
	}
	return nil
}

func writeByte(buf *bytes.Buffer, s rune) error {
	return binary.Write(buf, binary.BigEndian, int8(s))
}

func writeInt8(buf *bytes.Buffer, i int) error {
	return binary.Write(buf, binary.BigEndian, int8(i))
}

func writeInt32(buf *bytes.Buffer, i int) error {
	return binary.Write(buf, binary.BigEndian, int32(i))
}

func writeBytes(buf *bytes.Buffer, data []byte) error {
	return binary.Write(buf, binary.BigEndian, data)
}

func encodeString(buf *bytes.Buffer, s string) error {
	data := []byte(s)

	err := writeBytes(buf, []byte(strconv.Itoa(len(data))))
	if err != nil {
		return err
	}
	err = writeByte(buf, ':')
	if err != nil {
		return err
	}
	err = writeBytes(buf, data)
	if err != nil {
		return err
	}
	return nil
}
