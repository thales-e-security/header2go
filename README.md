# header2go

[![Build Status](https://travis-ci.com/thales-e-security/header2go.svg?branch=master)](https://travis-ci.com/thales-e-security/header2go)
[![Coverage Status](https://coveralls.io/repos/github/thales-e-security/header2go/badge.svg?branch=master)](https://coveralls.io/github/thales-e-security/header2go?branch=master)
[![GitHub](https://img.shields.io/github/license/mashape/apistatus.svg)](https://raw.githubusercontent.com/thales-e-security/header2go/master/LICENSE)


header2go generates skeleton Go implementations of C header files, which can be compiled with cgo to produce a shared
library.

## Installation

Please download a release from the releases page. We will release early and often, so building from source should
be unneeded.

If you want to build from source, you will first need to install dep. Then execute:

```
go get -d github.com/thales-e-security/header2go
cd $GOPATH/src/github.com/thales-e-security/header2go
dep ensure --vendor-only
go install .
```

## Usage

The command-line interface is very simple, just execute:

```
header2go <header-file> <output-dir>
```

The tool processes `<header-file>` and any included headers, before outputting the boilerplate code into `<output-dir>`.
The tool will happily overwrite existing files in `<output-dir>`, so be cautious.

If the processing completed successfully, it should be possible to compile the skeleton implementation using:

```
cd <output-dir>
go build -o library.so -buildmode=c-shared . 
```

which will produce a `library.so` file that can be linked into C programs.

## Documentation

Please see the [wiki](https://github.com/thales-e-security/header2go/wiki).

## Background

This project makes it easier to implement legacy APIs, described by C header files, in Go. The
original motivation was PKCS&nbsp;#11, a popular crypto API described by C header files and linked into dependent 
applications. We were interested in implementing this API using Go, but found the process of producing
the necessary CGo boiler-plate code rather tedious and error-prone.

header2go isn't yet capable of processing the whole PKCS&nbsp;#11 header set. 90% of the contents are successfully parsed at
present and the tool takes care to output as much useful boilerplate as possible, even in the face of errors.

### Output

header2go produces two files in the output directory. `main.go` contains all the CGo boilerplate, including the exported
functions and any necessary C types we have to include. The functions in this file handle the mapping between C memory
and Go memory in a Cgo compliant fashion. Each of the exported functions calls a corresponding Go function in `defs.go`.

The second file, `defs.go`, is the skeleton Go implementation of the API. The idea is to completely hide the horrors
of CGo in `main.go`, leaving `defs.go` as a pure Go implementation. A Go struct is created for each C type used in
the exported functions, ensuring that all the function parameters in `defs.go` are Go types.

To see examples of generated output, look in the `translate/testdata` directory. The subdirectories contain examples
of an input header file and the corresponding `main.go` and `defs.go`. These files are used for unit testing.

Here is a quick example to save you jumping straight into the source. For this simple header file:

```c
typedef struct {
    int val3;
} StructB;

typedef struct _StructA {
    int val1;
    StructB val2;
} StructA;

void functionA(StructA a1, long a2);
```

we get two files:

```go
// main.go

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

```

```go
// defs.go

package main

type StructB struct {
  Val3 int16
}

type StructA struct {
  Val1 int16
  Val2 StructB
}

func goFunctionA(a1 StructA, a2 int32) {
  // TODO implement goFunctionA
}

```

## Contributing

Contributions are very welcome. The tool is slowly growing in capability and there are plenty of tasks still to be done.
Before submitting a pull request, please ensure you have created a new test case (in `translate/testdata`) that tests
the generated code you expect to produce. The code in that directory should compile using cgo. Bonus points are awarded
if you include some C code that proves your generated code works, see  
`translate/testdata/08_pointers_to_structs_with_basic_types` for an example of that.

### TODO list (not exhaustive)

- [x] Command line interface
- [x] Basic function return types (no structs)
- [x] Pointer types as input parameters (copying into Go memory and back)
- [x] Fixed array types (e.g. char[16])
- [x] (Testing with) additional header files, i.e. #include statements
- [x] Void pointer parameter types (should map to unsafe.Pointer)
- [ ] Pointers to structs that contain pointers. (Makes the `convertFromXXX` functions more complicated).
- [ ] Function pointers (e.g. `typedef CK_RV(*CK_C_Finalize) (CK_VOID_PTR pReserved);`, which is a type called `CK_C_Finalize` that is a function returning `CK_RV` and taking an arguments of type `CK_VOID_PTR`.)
- [ ] Arbitrary return types
- [ ] Submit a PR to the github.com/elliotchance/c2go to refactor their code, ensuring we don't need to copy bits across to this project.

## Acknowledgements

This project uses the Clang AST parsing capabilities from https://github.com/elliotchance/c2go. Mostly this is achieved
by importing packages, however a small amount of code has been copied into `translate/astparse/parse.go`. The license
and copyright information is noted in that file.

