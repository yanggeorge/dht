package test

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"unicode"
)

type BDecode struct {
	arr []byte
	n   int
	i   int
}

func NewBDecode(arr []byte) BDecode {
	return BDecode{arr, len(arr), 0}
}

func (b *BDecode) Parse() (map[string]interface{}, error) {
	return b.Dic()
}

func (b *BDecode) Dic() (map[string]interface{}, error) {
	c, err := b.peek(0)
	if err != nil || c != 'd' {
		return nil, fmt.Errorf("error")
	}
	b.next()
	dic := make(map[string]interface{})
	var key string
	var val interface{}
	c, err = b.peek(0)
	if err != nil {
		return nil, fmt.Errorf("error")
	}
	var infoStart, infoEnd int
	for c != 'e' {
		key, err = b.string()
		if err != nil {
			return nil, fmt.Errorf("error")
		}

		if key == "info" {
			infoStart = b.i
		}

		val, err = b.element()
		if err != nil {
			return nil, fmt.Errorf("error")
		}

		if key == "info" {
			infoEnd = b.i
			infoHash := encodeToString(sha1.Sum(b.arr[infoStart:infoEnd]))
			dic["info_hash"] = infoHash
		}

		dic[key] = val
		c, err = b.peek(0)
		if err != nil {
			return nil, fmt.Errorf("error")
		}
	}

	c, err = b.next()
	if err != nil && c != 'e' {
		return nil, fmt.Errorf("error")
	}

	return dic, nil
}

func (b *BDecode) peek(n int) (rune, error) {
	if b.i+n < b.n {
		return rune(b.arr[b.i+n]), nil
	}
	return 0, errors.New("out of range")
}

func (b *BDecode) next() (rune, error) {
	if b.i < b.n {
		c := rune(b.arr[b.i])
		b.i += 1
		return c, nil
	}
	return 0, errors.New("out of range")
}

func (b *BDecode) string() (string, error) {
	num, err := b.num()
	if err != nil {
		return "", fmt.Errorf("error")
	}
	c, err := b.next()
	if err != nil || c != ':' {
		return "", fmt.Errorf("error")
	}

	s := string(b.arr[b.i : b.i+num])
	b.i += num
	return s, nil
}

func (b *BDecode) num() (int, error) {
	c, err := b.peek(0)
	if err != nil {
		return 0, fmt.Errorf("error")
	}

	num := 0
	for unicode.IsDigit(c) {
		num = num*10 + int(c) - '0'
		b.next()
		c, err = b.peek(0)
		if err != nil {
			return 0, fmt.Errorf("error")
		}
	}

	return num, nil
}

func (b *BDecode) element() (interface{}, error) {
	c, err := b.peek(0)
	if err != nil {
		return nil, fmt.Errorf("error")
	}

	if c == 'i' {
		return b.integer()
	} else if c == 'd' {
		return b.Dic()
	} else if c == 'l' {
		return b.list()
	} else {
		if unicode.IsDigit(c) {
			return b.string()
		}
		return nil, fmt.Errorf("error")
	}
}

func (b *BDecode) integer() (int, error) {
	c, err := b.next()
	if err != nil || c != 'i' {
		return 0, err
	}
	num, err := b.num()
	if err != nil {
		return 0, fmt.Errorf("error")
	}

	c, err = b.next()
	if err != nil || c != 'e' {
		return 0, fmt.Errorf("error")
	}
	return num, nil
}

func (b *BDecode) list() ([]interface{}, error) {
	var lis []interface{}

	c, err := b.next()
	if err != nil || c != 'l' {
		return lis, fmt.Errorf("error")
	}

	c, err = b.peek(0)
	if err != nil {
		return lis, err
	}
	for c != 'e' {
		ele, err := b.element()
		if err != nil {
			return lis, err
		}
		lis = append(lis, ele)
		c, err = b.peek(0)
		if err != nil {
			return lis, err
		}
	}
	b.next()
	return lis, nil
}
