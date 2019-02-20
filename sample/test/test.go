package main

import "fmt"

type A struct {
	a string
	b int
}

func (A) print() {

}

func main() {
	obj := &A{"abc", 1}
	obj2 := obj
	fmt.Printf("%v\n", obj.a)
	obj.a = "def"
	fmt.Printf("obj2=%v\n", obj2.a)
}
