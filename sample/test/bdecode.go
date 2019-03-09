package test

import (
	"container/list"
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

func (b BDecode) Parse() (map[string]interface{}, error) {
	return b.Dic()
}

func (b BDecode) Dic() (map[string]interface{}, error) {
	c, err := b.peek(0)
	if err != nil || c != 'd' {
		return nil, fmt.Errorf("error")
	}
	b.next()
	dic := make(map[string]interface{})
	var key string
	var val interface{}
	c, err = b.peek(0)
	for err != nil && c != 'e' {
		key, err = b.string()
		if err != nil {
			return nil, fmt.Errorf("error")
		}
		val, err = b.element()
		if err != nil {
			return nil, fmt.Errorf("error")
		}
		dic[key] = val
		c, err = b.peek(0)
	}

	c, err = b.next()
	if err != nil && c != 'e' {
		return nil, fmt.Errorf("error")
	}

	return dic, nil
}

func (b BDecode) peek(n int) (rune, error) {
	if b.i+n+1 < b.n {
		return rune(b.arr[b.i+n]), nil
	}
	return 0, errors.New("out of range")
}

func (b BDecode) next() (rune, error) {
	if b.i+1 < b.n {
		c := rune(b.arr[b.i])
		b.i += 1
		return c, nil
	}
	return 0, errors.New("out of range")
}

func (b BDecode) string() (string, error) {
	num, err := b.num()
	if err != nil {
		return "", fmt.Errorf("error")
	}
	c, err := b.next()
	if err != nil || c != ':' {
		return "", fmt.Errorf("error")
	}

	s := string(b.arr[b.i : b.i+num])

	return s, nil
}

func (b BDecode) num() (int, error) {
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

func (b BDecode) element() (interface{}, error) {
	c, err := b.peek(0)
	if err != nil {
		return nil, fmt.Errorf("error")
	}

	switch c {
	case 'i':
		return b.integer()
	case 'd':
		return b.Dic()
	case 'l':
		return b.list()
	case '0':
	case '1':
	case '2':
	case '3':
	case '4':
	case '5':
	case '6':
	case '7':
	case '8':
	case '9':
		return b.string()
	default:
		return nil, fmt.Errorf("error")
	}

	return nil, fmt.Errorf("error")
}

func (b BDecode) integer() (int, error) {
	c, err := b.next()
	if err != nil || c != 'i' {
		return 0, fmt.Errorf("error")
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

func (b BDecode) list() (*list.List, error) {
	lis := list.New()

	c, err := b.next()
	if err != nil || c != 'l' {
		return lis, fmt.Errorf("error")
	}

	c, err = b.peek(0)
	if err != nil {
		return lis, fmt.Errorf("error")
	}
	for c != 'e' {
		ele, err := b.element()
		if err != nil {
			return lis, fmt.Errorf("error")
		}
		lis.PushBack(ele)
	}
	return lis, nil
}
