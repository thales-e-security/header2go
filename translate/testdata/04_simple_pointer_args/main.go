package main

/*

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

//export functionA
func functionA(a1 *C.int, a2 C.long) C.int {
	var a1Array *Int16Array
	if a1 != nil {
		a1Array = &Int16Array{cdata: a1}
	}

	res := goFunctionA(a1Array, int32(a2))
	a1Array.writeBack()
	return C.int(res)
}

func main() {}
