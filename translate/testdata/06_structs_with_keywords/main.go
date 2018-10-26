package main

/*

typedef struct _StructA {
  int type;
} StructA;

*/
import "C"

func convertToStructA(s C.StructA) StructA {
	return StructA{
		Type: int16(s._type),
	}
}

//export functionA
func functionA(a1 C.StructA, a2 C.long) {
	goFunctionA(convertToStructA(a1), int32(a2))
}

func main() {}
