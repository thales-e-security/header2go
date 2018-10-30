package main

type StructB struct {
	Val1  byte
	Val2  int8
	Val3  byte
	Val4  int16
	Val5  uint16
	Val6  int16
	Val7  uint16
	Val8  int32
	Val9  uint32
	Val10 int64
	Val11 uint64
	Val12 float32
	Val13 float64
}

type StructA struct {
	Val1 int16
	Val2 StructB
}

func goFunctionA(a1 StructA, a2 int32) {
	// TODO implement goFunctionA
}
