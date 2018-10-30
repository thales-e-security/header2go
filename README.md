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
header2go <header-file> <output-dir> [<config-file>]
```

The tool processes `<header-file>` and any included headers, before outputting the boilerplate code into `<output-dir>`.
The tool will happily overwrite existing files in `<output-dir>`, so be cautious.

Void pointers are handled by mappings in a configuration file. See [the wiki](https://github.com/thales-e-security/header2go/wiki/Void-Pointers).

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
- [ ] Generate a suitable config file, based on scanning for void pointers.

## Acknowledgements

This project uses the Clang AST parsing capabilities from https://github.com/elliotchance/c2go. Mostly this is achieved
by importing packages, however a small amount of code has been copied into `translate/astparse/parse.go`. The license
and copyright information is noted in that file.

