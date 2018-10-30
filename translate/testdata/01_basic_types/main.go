package main

/*

 */
import "C"

//export functionA
func functionA(a1 C.char, a2 C.schar, a3 C.uchar, a4 C.short, a5 C.ushort, a6 C.int, a7 C.uint, a8 C.long, a9 C.ulong, a10 C.longlong, a11 C.ulonglong, a12 C.float, a13 C.double) {
	goFunctionA(byte(a1), int8(a2), byte(a3), int16(a4), uint16(a5), int16(a6), uint16(a7), int32(a8), uint32(a9), int64(a10), uint64(a11), float32(a12), float64(a13))
}

func main() {}
