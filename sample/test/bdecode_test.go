package test

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func TestString(t *testing.T) {
	bs := []byte("5:hello")
	b := NewBDecode(bs)
	s, err := b.string()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s)
}

func TestInteger(t *testing.T) {
	bs := []byte("i343e")
	b := NewBDecode(bs)
	i, err := b.integer()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(i)
}

func TestParse1(t *testing.T) {
	bs := []byte("d4:spaml1:a1:bee")
	b := NewBDecode(bs)
	dic, e := b.Parse()
	check(e)
	fmt.Println(dic)
}

func TestParse(t *testing.T) {
	path := "/Users/ym/tmp/venom.torrent"
	dat, e := ioutil.ReadFile(path)
	check(e)
	b := NewBDecode(dat)
	dic, e := b.Parse()
	check(e)
	fmt.Println(dic)
	val := dic["info"]
	infoDic := val.(map[string]interface{})
	fmt.Println(infoDic["pieces"])
}
