package test

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	dic := make(map[string]interface{}, 10)
	dic2 := make(map[string]interface{}, 10)
	dic["abc"] = 1
	dic["bcd"] = "ym"
	lis := make([]interface{}, 0)
	lis = append(lis, 1)
	lis = append(lis, 2)
	dic2["a"] = "abc"
	lis = append(lis, dic2)
	dic["list"] = lis

	if reflect.TypeOf(lis) == reflect.TypeOf([]interface{}(nil)) {
		fmt.Println("yes")
	}

	data, err := encodeToBin(dic)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range data {
		fmt.Printf("%d,", v)
	}
	fmt.Printf("\n")
	b := NewBDecode(data)
	dic, err = b.Dic()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dic)
}
