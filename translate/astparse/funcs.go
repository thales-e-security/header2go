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

package astparse

import (
	"log"
	"regexp"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/thales-e-security/header2go/translate/errors"
)

// CFuncDef describes a function found in the C header file.
type CFuncDef struct {

	// Name is the name of the function in the header file.
	Name string

	// Params defines the parameters and their type.
	Params []*CStructInstance

	// ReturnType is the function return type, or nil for void functions.
	ReturnType *CStructInstance
}

// ProcessFuncs parses the function list and builds a set of function descriptions. It also figures out the
// minimal used set of structs and returns this.
func ProcessFuncs(context *errors.ParseContext, functions []*ast.FunctionDecl,
	structs []*CStruct, types []*CType) ([]*CFuncDef, []*CStruct) {

	// Create a map of structs by name
	structsByName := make(map[string]*CStruct)
	for _, s := range structs {
		// Some structs don't have a name, but that's fine - they can't be referenced directly in a function definition.
		if s.Name != "" {
			structsByName["struct "+s.Name] = s
		}
	}

	// Create a map of types by name
	typesByName := make(map[string]*CType)
	for _, t := range types {
		typesByName[t.Name] = t
	}

	var funcs []*CFuncDef
	usedStructs := make(map[*CStruct]bool)

functionLoop:
	for _, funcDecl := range functions {

		f := &CFuncDef{
			Name: funcDecl.Name,
		}

		// Process params
		for _, n := range funcDecl.ChildNodes {
			if paramDecl, ok := n.(*ast.ParmVarDecl); ok {

				p := findType(paramDecl.Type, structsByName, typesByName)
				if p == nil {
					context.AddFunctionError(funcDecl.Pos,
						"cannot process function %s: cannot resolve type '%s' used in param %s", funcDecl.Name,
						paramDecl.Type, paramDecl.Name)
					continue functionLoop
				}

				// Recursively flag structs as being used
				var markStructAsSeen func(s *CStruct)

				markStructAsSeen = func(s *CStruct) {
					if usedStructs[s] {
						return
					}

					for _, f := range s.Fields {
						markStructAsSeen(f.Struct)
					}

					if !s.Basic {
						usedStructs[s] = true
					}
				}
				markStructAsSeen(p.Struct)

				p.Name = paramDecl.Name
				f.Params = append(f.Params, p)
			} else {
				log.Printf("Skipping child node of type %T", n)
			}
		}

		// Process return type
		funcType := parseReturnType(funcDecl.Type)
		if funcType == "" {
			context.AddFunctionError(funcDecl.Pos,
				"cannot process function %s: cannot parse function type '%s'", funcDecl.Name,
				funcDecl.Type)
			continue
		}

		if funcType != "void" {
			f.ReturnType = findType(funcType, structsByName, typesByName)
			if f.ReturnType == nil {
				context.AddFunctionError(funcDecl.Pos,
					"cannot process function %s: cannot resolve function return type '%s'", funcDecl.Name,
					funcType)
				continue
			}

			// Note that we've seen this type
			if !f.ReturnType.Struct.Basic {
				usedStructs[f.ReturnType.Struct] = true
			}
		}
		funcs = append(funcs, f)
	}

	// Preserve type order
	var usedStructList []*CStruct
	for _, s := range structs {
		if usedStructs[s] {
			usedStructList = append(usedStructList, s)
		}
	}

	return funcs, usedStructList
}

func parseReturnType(functionType string) string {
	// Function type is of the form "<return type> (<param 1 type>, ...)"
	re := regexp.MustCompile(`^(.*)\(.*\)$`)
	submatches := re.FindStringSubmatch(functionType)
	if submatches == nil {
		return ""
	}

	return strings.TrimSpace(submatches[1]) // [0] is the whole match
}

func findType(typeName string, structsByName map[string]*CStruct,
	typesByName map[string]*CType) *CStructInstance {

	p := &CStructInstance{}

	// Count pointers
	nonPointerType, numPointers := stripPointers(typeName)
	p.PointerCount = numPointers

	p.Struct = processBasic(nonPointerType)
	if p.Struct != nil {
		return p
	}

	p.Struct = structsByName[nonPointerType]
	if p.Struct != nil {
		return p
	}

	t := typesByName[nonPointerType]
	if t == nil {
		// We have failed to find this type
		return nil
	}

	p.PointerCount += t.PointerCount
	p.Struct = t.Struct
	return p
}
