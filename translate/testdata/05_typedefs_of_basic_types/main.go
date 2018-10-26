package main

/*

 */
import "C"

//export functionA
func functionA(a1 C.int) {
	goFunctionA(int16(a1))
}

func main() {}
