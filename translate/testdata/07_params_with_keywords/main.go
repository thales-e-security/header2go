package main

/*

 */
import "C"

//export functionA
func functionA(_type C.int, a2 C.long) {
	goFunctionA(int16(_type), int32(a2))
}

func main() {}
