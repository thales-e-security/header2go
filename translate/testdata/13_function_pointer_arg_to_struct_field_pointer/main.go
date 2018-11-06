package main

/*

typedef struct  {
  int *buffer;
  int len;
} Buffer;

*/
import "C"
import "unsafe"

type BufferArray struct {
	cdata  *C.Buffer
	godata []Buffer
}

func (a *BufferArray) Slice(len int) []Buffer {
	if a.godata != nil {
		return a.godata
	}

	a.godata = make([]Buffer, len)
	for i := 0; i < len; i++ {
		a.godata[i] = convertToBuffer(*(*C.Buffer)(unsafe.Pointer(uintptr(unsafe.Pointer(a.cdata)) + uintptr(i)*unsafe.Sizeof(*a.cdata))))
	}

	return a.godata
}

func (a *BufferArray) writeBack() {
	for i := range a.godata {
		*(*C.Buffer)(unsafe.Pointer(uintptr(unsafe.Pointer(a.cdata)) + uintptr(i)*unsafe.Sizeof(*a.cdata))) = convertFromBuffer(a.godata[i])
	}
}

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

func convertFromBuffer(s Buffer) C.Buffer {
	s.Buffer.writeBack()
	return C.Buffer{
		buffer: s.Buffer.cdata,
		len:    C.int(s.Len),
	}
}

//export functionA
func functionA(buffer *C.Buffer) C.int {
	var bufferArray *BufferArray
	if buffer != nil {
		bufferArray = &BufferArray{cdata: buffer}
	}

	res := goFunctionA(bufferArray)
	bufferArray.writeBack()
	return C.int(res)
}

func main() {}
