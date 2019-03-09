package test

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	bs := make([]byte, 12)
	bs = []byte("hello")
	b := NewBDecode(bs)
	c, err := b.peek(12)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c)
	a := '0'
	fmt.Println(int(a) - '0')
}
