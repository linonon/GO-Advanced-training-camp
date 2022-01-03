package main

import "fmt"

type IceCreamMaker interface {
	Hello()
}

type Ben struct {
	id   int
	name string
}

func (b *Ben) Hello() {
	fmt.Printf("Ben says, \"Hello my name is %s\"\n", b.name)
}

// 這種情況的 Jerry definition 也不會報錯。
// type Jerry struct {
// 	field1 *[5]byte
// 	field1 int
// }

type Jerry struct {
	name string
}

func (j *Jerry) Hello() {
	fmt.Printf("Jerry Says, \"Hello my name is %s\"\n", j.name)
}

func main() {
	var ben = &Ben{id: 10, name: "Ben"}
	var jerry = &Jerry{"Jerry"}
	var maker IceCreamMaker = ben

	var loop0, loop1 func()

	loop0 = func() {
		maker = ben
		go loop1()
	}

	loop1 = func() {
		maker = jerry
		go loop0()
	}

	go loop0()

	for {
		maker.Hello()
	}
}
