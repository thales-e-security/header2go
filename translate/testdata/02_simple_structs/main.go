package main

/*

typedef struct  {
  char val1;
  signed char val2;
  unsigned char val3;
  short val4;
  unsigned short val5;
  int val6;
  unsigned int val7;
  long val8;
  unsigned long val9;
  long long val10;
  unsigned long long val11;
  float val12;
  double val13;
} StructB;

typedef struct _StructA {
  int val1;
  StructB val2;
} StructA;

*/
import "C"

func convertToStructB(s C.StructB) StructB {
	return StructB{
		Val1:  byte(s.val1),
		Val2:  int8(s.val2),
		Val3:  byte(s.val3),
		Val4:  int16(s.val4),
		Val5:  uint16(s.val5),
		Val6:  int16(s.val6),
		Val7:  uint16(s.val7),
		Val8:  int32(s.val8),
		Val9:  uint32(s.val9),
		Val10: int64(s.val10),
		Val11: uint64(s.val11),
		Val12: float32(s.val12),
		Val13: float64(s.val13),
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
