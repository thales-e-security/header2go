package main

/*

typedef struct  {
  int *buffer;
  int len;
} Buffer;

*/
import "C"
import "unsafe"

type Int16Array struct {
	cdata  *C.int
	godata []int16
}

func (a *Int16Array) Slice(len int) []int16 {
	if a.godata != nil {
		return a.godata
	}

	a.godata = make([]int16, len)
	for i := 0; i < len; i++ {
		a.godata[i] = int16(*(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(a.cdata)) + uintptr(i)*unsafe.Sizeof(*a.cdata))))
	}

	return a.godata
}

func (a *Int16Array) writeBack() {
	for i := range a.godata {
		*(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(a.cdata)) + uintptr(i)*unsafe.Sizeof(*a.cdata))) = C.int(a.godata[i])
	}
}

func convertToBuffer(s C.Buffer) Buffer {
	return Buffer{
		Buffer: Int16Array{cdata: s.buffer},
		Len:    int16(s.len),
	}
}

//export functionA
func functionA(buffer C.Buffer) C.int {
	res := goFunctionA(convertToBuffer(buffer))
	return C.int(res)
}

func main() {}
