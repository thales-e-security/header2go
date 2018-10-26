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

package translate

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"log"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/elliotchance/c2go/ast"
	"github.com/pkg/errors"
	"github.com/thales-e-security/header2go/translate/astparse"
	errs "github.com/thales-e-security/header2go/translate/errors"
)

// templateData contains the data needed to populate the source templates.
type templateData struct {
	ArrayTypes             []templateArrayType
	TypesNeedingConversion []*templateType
	PointerTypes           []*templateType
	Types                  []*templateType
	Funcs                  []*templateFunc
}

// Translate processes the AST data (function definitions, type definitions) to produce the output files.
func Translate(context *errs.ParseContext, functions []*ast.FunctionDecl, recordDeclarations []*ast.RecordDecl,
	typeDeclarations []*ast.TypedefDecl, outputDir string) error {

	// Process the raw AST data and convert into our representation of structs and types
	structs, types := astparse.ProcessTypes(context, recordDeclarations, typeDeclarations)

	// Process the raw AST data for functions, noting which structs are actually used (in the signatures).
	funcs, usedStructs := astparse.ProcessFuncs(context, functions, structs, types)

	// Find the types used as pointers (and also log an error if we find double-pointers or return-type pointers)
	pointerTypes, typesNeedingConversion, funcs := listPointerTypes(context, funcs)

	arrayTypes := listArrayTypes(context, usedStructs)

	var templateTypes []*templateType
	for _, s := range usedStructs {
		tt := templateType(*s)
		templateTypes = append(templateTypes, &tt)
	}

	var templateFuncs []*templateFunc
	for _, f := range funcs {
		tf := templateFunc(*f)
		templateFuncs = append(templateFuncs, &tf)
	}

	data := templateData{
		ArrayTypes:             arrayTypes,
		TypesNeedingConversion: typesNeedingConversion,
		PointerTypes:           pointerTypes,
		Types:                  templateTypes,
		Funcs:                  templateFuncs,
	}

	err := writeTemplateToFile("main.gotmpl", path.Join(outputDir, "main.go"), data)
	if err != nil {
		return errors.Wrap(err, "failed to write main.go template")
	}

	err = writeTemplateToFile("defs.gotmpl", path.Join(outputDir, "defs.go"), data)
	return errors.Wrap(err, "failed to write defs.go template")
}

// listArrayTypes parses the types and finds each array type + size combination.
func listArrayTypes(context *errs.ParseContext, types []*astparse.CStruct) (arrayTypes []templateArrayType) {
	foundTypes := make(map[*astparse.CStruct]map[uint]bool)

	for _, t := range types {
		for _, f := range t.Fields {
			if f.ArrayCount > 0 {

				if foundTypes[f.Struct] == nil {
					foundTypes[f.Struct] = make(map[uint]bool)
				}
				foundTypes[f.Struct][f.ArrayCount] = true
			}
		}
	}

	for t := range foundTypes {
		for c := range foundTypes[t] {
			tt := templateType(*t)
			arrayTypes = append(arrayTypes, templateArrayType{ArrayCount: c, Type: &tt})
		}
	}

	// We want a predictable order, so let's alphabetically sort
	sort.Slice(arrayTypes, func(i, j int) bool {
		return strings.Compare(arrayTypes[i].Type.CanonicalName(), arrayTypes[j].Type.CanonicalName()) <= 0
	})
	return
}

// listPointerTypes parses the function definitions and looks for parameters that use a pointer type.
func listPointerTypes(context *errs.ParseContext, funcs []*astparse.CFuncDef) (pointerTypes,
	pointerFields []*templateType, updatedFuncs []*astparse.CFuncDef) {

	updatedFuncs = funcs
	pointerTypeMap := make(map[*astparse.CStruct]bool)

	// We do a downwards for loop, so we can remove functions that cause errors as we iterate
outerLoop:
	for i := len(funcs) - 1; i >= 0; i-- {
		if funcs[i].ReturnType != nil && funcs[i].ReturnType.PointerCount != 0 {
			context.AddFunctionError(ast.Position{}, "Cannot support pointer return type for function %s", funcs[i].Name)

			// Remove this function
			updatedFuncs = append(updatedFuncs[:i], updatedFuncs[i+1:]...)
			continue
		}

		for _, p := range funcs[i].Params {
			if p.PointerCount == 1 {
				pointerTypeMap[p.Struct] = true
			} else if p.PointerCount > 1 {
				context.AddFunctionError(ast.Position{}, "Cannot support double pointer parameter type for function %s", funcs[i].Name)

				// Remove this function
				updatedFuncs = append(updatedFuncs[:i], updatedFuncs[i+1:]...)
				continue outerLoop
			}
		}
	}

	var rawPointerTypes []*astparse.CStruct
	for s := range pointerTypeMap {
		rawPointerTypes = append(rawPointerTypes, s)
		tt := templateType(*s)
		pointerTypes = append(pointerTypes, &tt)
	}
	// We want a predictable order, so let's alphabetically sort
	sort.Slice(pointerTypes, func(i, j int) bool {
		return strings.Compare(pointerTypes[i].CanonicalName(), pointerTypes[j].CanonicalName()) <= 0
	})

	// Now parse the (raw) pointer types, looking for functions that will need a "convertFromXXX" function. Each pointer
	// type results in a <Foo>Array type being created and each field of that type will require a conversion
	// function (so we can convert back from the Go struct equivalent to the C struct).
	rawPointerFields := listPointerFields(context, rawPointerTypes)
	for _, v := range rawPointerFields {
		tt := templateType(*v)
		pointerFields = append(pointerFields, &tt)
	}

	// We want a predictable order, so let's alphabetically sort
	sort.Slice(pointerFields, func(i, j int) bool {
		return strings.Compare(pointerFields[i].CanonicalName(), pointerFields[j].CanonicalName()) <= 0
	})

	return
}

// listPointerFields recursively examines each pointer type and creates a list of all the types referenced
// as struct fields.
func listPointerFields(context *errs.ParseContext, pointerTypes []*astparse.CStruct) []*astparse.CStruct {

	foundTypes := make(map[*astparse.CStruct]bool)

	var f func(*astparse.CStruct)

	f = func(s *astparse.CStruct) {

		if s.Basic {
			// We don't need to add conversion functions for basic types
			return
		}

		// Mark this type as needing a conversion function
		foundTypes[s] = true

		for _, field := range s.Fields {
			if field.PointerCount > 0 {
				context.AddTypeError(ast.Position{}, "Cannot create conversion function for type %s due to "+
					"having a pointer field %s.", s.Name, field.Name)
				continue
			}

			// Recursively process sub fields
			f(field.Struct)
		}
	}

	for _, t := range pointerTypes {
		f(t)
	}

	var res []*astparse.CStruct
	for t := range foundTypes {
		res = append(res, t)
	}

	return res
}

func writeTemplateToFile(templateName, outputFile string, data interface{}) error {

	outputBytes, err := writeTemplate(templateName, data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputFile, outputBytes, 0644)
}

func writeTemplate(templateName string, data interface{}) ([]byte, error) {
	templateData, err := Asset(path.Join("translate", "templates", templateName))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read template")
	}

	tmpl, err := template.New("").Parse(string(templateData))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse template")
	}

	buffer := new(bytes.Buffer)

	err = tmpl.Execute(buffer, data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute template")
	}

	formatted, err := format.Source(buffer.Bytes())

	if err != nil {
		log.Println("WARNING: failed to format source. This means the generation has probably gone wrong somewhere.")
		log.Println(err.Error())
		return buffer.Bytes(), nil
	}

	return formatted, nil
}

func capitalise(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func basicToGo(ctype string) string {
	switch ctype {
	case "char":
		return "byte"

	case "signed char":
		return "byte"

	case "unsigned char":
		return "int8"

	case "short", "short int", "signed short", "signed short int":
		return "int16"

	case "unsigned short", "unsigned short int":
		return "uint16"

	case "int", "signed", "signed int":
		return "int16"

	case "unsigned", "unsigned int":
		return "uint16"

	case "long", "long int", "signed long", "signed long int":
		return "int32"

	case "unsigned long", "unsigned long int":
		return "uint32"

	case "long long", "long long int", "signed long long", "signed long long int":
		return "int64"

	case "unsigned long long", "unsigned long long int":
		return "uint64"

	case "float":
		return "float32"

	case "double":
		return "float64"

	default:
		// should never happen
		return "ERROR"
	}
}
