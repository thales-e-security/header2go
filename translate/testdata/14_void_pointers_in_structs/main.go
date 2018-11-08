package main

/*

typedef struct  {
  unsigned long len;
  void *ptr;
} TypeA;

*/
import "C"
import "unsafe"

type ByteArray struct {
	cdata  *C.uchar
	godata []byte
}

func (a *ByteArray) Slice(len int) []byte {
	if a.godata != nil {
		return a.godata
	}

	a.godata = make([]byte, len)
	for i := 0; i < len; i++ {
		a.godata[i] = byte(*(*C.uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(a.cdata)) + uintptr(i)*unsafe.Sizeof(*a.cdata))))
	}

	return a.godata
}

func (a *ByteArray) writeBack() {
	for i := range a.godata {
		*(*C.uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(a.cdata)) + uintptr(i)*unsafe.Sizeof(*a.cdata))) = C.uchar(a.godata[i])
	}
}

func convertToTypeA(s C.TypeA) TypeA {
	return TypeA{
		Len: uint32(s.len),
		Ptr: ByteArray{cdata: (*C.uchar)(s.ptr)},
	}
}

//export functionA
func functionA(a1 C.TypeA, a2 C.long) C.int {
	res := goFunctionA(convertToTypeA(a1), int32(a2))
	return C.int(res)
}

func main() {}
