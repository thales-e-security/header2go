package main

/*

typedef struct _StructA {
  int val1;
  int val2[10];
} StructA;

*/
import "C"
import "unsafe"

type StructAArray struct {
	cdata  *C.StructA
	godata []StructA
}

func (a *StructAArray) Slice(len int) []StructA {
	if a.godata != nil {
		return a.godata
	}

	a.godata = make([]StructA, len)
	for i := 0; i < len; i++ {
		a.godata[i] = convertToStructA(*(*C.StructA)(unsafe.Pointer(uintptr(unsafe.Pointer(a.cdata)) + uintptr(i)*unsafe.Sizeof(*a.cdata))))
	}

	return a.godata
}

func (a *StructAArray) writeBack() {
	for i := range a.godata {
		*(*C.StructA)(unsafe.Pointer(uintptr(unsafe.Pointer(a.cdata)) + uintptr(i)*unsafe.Sizeof(*a.cdata))) = convertFromStructA(a.godata[i])
	}
}

func convertToStructA(s C.StructA) StructA {
	return StructA{
		Val1: int16(s.val1),
		Val2: convertToint16Array10(s.val2),
	}
}

func convertFromStructA(s StructA) C.StructA {
	return C.StructA{
		val1: C.int(s.Val1),
		val2: [10]C.int{C.int(s.Val2[0]), C.int(s.Val2[1]), C.int(s.Val2[2]), C.int(s.Val2[3]), C.int(s.Val2[4]), C.int(s.Val2[5]), C.int(s.Val2[6]), C.int(s.Val2[7]), C.int(s.Val2[8]), C.int(s.Val2[9])},
	}
}

func convertToint16Array10(a [10]C.int) (res [10]int16) {
	for i := range a {
		res[i] = int16(a[i])
	}
	return
}

//export functionA
func functionA(a1 *C.StructA, a2 C.long) {
	var a1Array *StructAArray
	if a1 != nil {
		a1Array = &StructAArray{cdata: a1}
	}

	goFunctionA(a1Array, int32(a2))
}

func main() {}
