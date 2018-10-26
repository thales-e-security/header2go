package main

/*

typedef struct  {
  int val3[10];
} StructB;

typedef struct _StructA {
  int val1;
  StructB val2[2];
} StructA;

*/
import "C"

func convertToStructB(s C.StructB) StructB {
	return StructB{
		Val3: convertToint16Array10(s.val3),
	}
}

func convertToStructA(s C.StructA) StructA {
	return StructA{
		Val1: int16(s.val1),
		Val2: convertToStructBArray2(s.val2),
	}
}

func convertToStructBArray2(a [2]C.StructB) (res [2]StructB) {
	for i := range a {
		res[i] = convertToStructB(a[i])
	}
	return
}

func convertToint16Array10(a [10]C.int) (res [10]int16) {
	for i := range a {
		res[i] = int16(a[i])
	}
	return
}

//export functionA
func functionA(a1 C.StructA, a2 C.long) {
	goFunctionA(convertToStructA(a1), int32(a2))
}

func main() {}
