package main

/*

typedef struct  {
  int a;
} TypeA;

*/
import "C"
import "unsafe"

type TypeAArray struct {
	cdata  *C.TypeA
	godata []TypeA
}

func (a *TypeAArray) Slice(len int) []TypeA {
	if a.godata != nil {
		return a.godata
	}

	a.godata = make([]TypeA, len)
	for i := 0; i < len; i++ {
		a.godata[i] = convertToTypeA(*(*C.TypeA)(unsafe.Pointer(uintptr(unsafe.Pointer(a.cdata)) + uintptr(i)*unsafe.Sizeof(*a.cdata))))
	}

	return a.godata
}

func (a *TypeAArray) writeBack() {
	for i := range a.godata {
		*(*C.TypeA)(unsafe.Pointer(uintptr(unsafe.Pointer(a.cdata)) + uintptr(i)*unsafe.Sizeof(*a.cdata))) = convertFromTypeA(a.godata[i])
	}
}

func convertToTypeA(s C.TypeA) TypeA {
	return TypeA{
		A: int16(s.a),
	}
}

func convertFromTypeA(s TypeA) C.TypeA {
	return C.TypeA{
		a: C.int(s.A),
	}
}

//export functionA
func functionA(a1 unsafe.Pointer, a2 C.long) C.int {
	var a1Array *TypeAArray
	if a1 != nil {
		a1Array = &TypeAArray{cdata: (*C.TypeA)(a1)}
	}

	res := goFunctionA(a1Array, int32(a2))
	a1Array.writeBack()
	return C.int(res)
}

func main() {}
