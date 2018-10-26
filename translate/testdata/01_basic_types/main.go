package main

/*

 */
import "C"

//export functionA
func functionA(a1 C.int, a2 C.long) {
	goFunctionA(int16(a1), int32(a2))
}

func main() {}
