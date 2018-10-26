// Copyright 2018 Thales UK Limited
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
// documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions
// of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"log"
	"os"

	"git.tesfabric.com/goethite/header2go/translate/errors"

	"git.tesfabric.com/goethite/header2go/translate"
)

//go:generate go-bindata -o translate/bindata.go -pkg translate translate/templates/...
func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: header2go <header-file> <output-dir>")
	}

	headerFile := os.Args[1]
	outputDir := os.Args[2]

	if _, err := os.Stat(outputDir); err == nil {
		log.Printf("Overwriting existing path '%s'", outputDir)
	} else {
		err := os.MkdirAll(outputDir, 0750)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	context := errors.NewParseContext()
	err := translate.Process(context, headerFile, outputDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	typeErrors := context.TypeErrors()
	if len(typeErrors) > 0 {
		log.Println("Type processing errors:")

		for _, e := range typeErrors {
			log.Println("\t" + e.Error())
		}
	}

	functionErrors := context.FunctionErrors()
	if len(functionErrors) > 0 {
		log.Println("Function processing errors:")

		for _, e := range functionErrors {
			log.Println("\t" + e.Error())
		}
	}

	log.Println("Done")
}
