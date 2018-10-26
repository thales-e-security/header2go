package main

/*

typedef struct  {
  int val3;
} StructB;

typedef struct _StructA {
  int val1;
  StructB val2;
} StructA;

*/
import "C"

func convertToStructB(s C.StructB) StructB {
	return StructB{
		Val3: int16(s.val3),
	}
}

func convertToStructA(s C.StructA) StructA {
	return StructA{
		Val1: int16(s.val1),
		Val2: convertToStructB(s.val2),
	}
}

//export functionA
func functionA(a1 C.StructA, a2 C.long) {
	goFunctionA(convertToStructA(a1), int32(a2))
}

func main() {}
