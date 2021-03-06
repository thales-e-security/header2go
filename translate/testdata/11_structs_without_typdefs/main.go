package main

/*

typedef struct  {
  int val3;
} StructB;

struct StructA {
  int val1;
  StructB val2;
} ;

*/
import "C"

func convertToStructB(s C.StructB) StructB {
	return StructB{
		Val3: int16(s.val3),
	}
}

func convertToStructA(s C.struct_StructA) StructA {
	return StructA{
		Val1: int16(s.val1),
		Val2: convertToStructB(s.val2),
	}
}

//export functionA
func functionA(a1 C.struct_StructA, a2 C.long) {
	goFunctionA(convertToStructA(a1), int32(a2))
}

func main() {}
