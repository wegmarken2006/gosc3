package main

import (
	"fmt"
	. "gosc3/osc"
	. "gosc3/sc3"
	"math"
)

func main() {

	fmt.Println("start")
	//	s1 := "gigi"
	//	fmt.Println(StrPstr(s1))

	a1 := DecodeF32(EncodeF32(51.33))
	fmt.Println(a1)
	//	ScStart()

	PrintUgen(UAbs(13))
	PrintUgen(UAdd(NewIConst(6), NewIConst(2)))
	fmt.Println("\nEnd")

	/*
		out := 2
	*/
}

func UAbs(ugen interface{}) UgenType {
	return MkUnaryOperator(5, math.Abs, ugen)
}
func UAdd(op1 interface{}, op2 interface{}) UgenType {
	fun := func(x float64, y float64) float64 { return x + y }
	return MkBinaryOperator(0, fun, op1, op2)
}
