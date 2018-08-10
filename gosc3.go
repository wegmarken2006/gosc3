package main

import (
	"fmt"
	. "gosc3/osc"
	//		. "gosc3/sc3"
)

func main() {

	fmt.Println("start")
	//	s1 := "gigi"
	//	fmt.Println(StrPstr(s1))

	a1 := DecodeF32(EncodeF32(51.33))
	fmt.Println(a1)
	OscSend([]byte{1, 2, 3})
	fmt.Println("end")

	/*
		out := 2
	*/
}
