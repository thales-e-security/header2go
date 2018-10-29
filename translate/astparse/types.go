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
	"strconv"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/thales-e-security/header2go/translate/errors"
)

// A CStructInstance contains information about a specific use of a struct, for instance as a field in another struct,
// as a parameter in a function or as a return type.
type CStructInstance struct {

	// Name contains the name of the variable or parameter. This will be empty in contexts where a name is not
	// associated with the struct (e.g. return types).
	Name string

	// PointerCount counts how many pointers were used in this instance. E.g. "int *foo" would result in a PointerCount
	// of 1.
	PointerCount uint

	// ArrayCount indicates if the struct was used in a fixed-sized array. E.g. "int foo[4]" would result in an
	// ArrayCount of 4, where as "int foo" would result in 0.
	ArrayCount uint

	// WasVoidPointer indicates if this instance was actually a void pointer, which we've mapped to an underlying
	// type through the user-supplied config.
	WasVoidPointer bool

	// Struct is the type that was used in this instance.
	Struct *CStruct
}

// A CType describes a typedef for a CStruct.
type CType struct {

	// Name is the typedef name.
	Name string

	// PointerCount describes the number of pointers this typedef introduces, on top of the underlying struct.
	PointerCount uint

	// Struct is the underlying struct this typedef references. The parsing phase will resolve any intermediate
	// typdefs, ensuring we can directly tie each CType to a CStruct.
	Struct *CStruct
}

// A CStruct represents a C struct.
type CStruct struct {

	// Basic indicates if this is a pseudo-struct representing a basic type (e.g. int, char etc.).
	Basic bool

	// Name contains the name of the struct, which may be empty for anonymous structs.
	Name string

	// TypeName contains the type name for this struct, if one exists. During processing, if there are any typedefs
	// which alias this struct (without pointers), one of them will be stored in this field.
	TypeName string

	// Fields contains the fields defined in this struct.
	Fields []*CStructInstance
}

// processTypeNode recursively examines a typedef tree structure and locates the address of the underlying
// struct. It also keeps track of the number of pointers found.
func processTypeNode(context *errors.ParseContext, n ast.Node, newType *CType) (address *ast.Address) {
	switch v := n.(type) {
	case *ast.PointerType:
		newType.PointerCount++
		return processTypeNode(context, v.Children()[0], newType)

	case *ast.TypedefType, *ast.ElaboratedType, *ast.RecordDecl, *ast.RecordType:
		// Recurse deeper
		for _, c := range v.Children() {
			mark := context.Mark()
			address = processTypeNode(context, c, newType)
			if address != nil || context.HasErrors(mark) {
				return address
			}
		}

		// We shouldn't get here without having found something (even a basic type)...
		if newType.Struct == nil {
			context.AddTypeError(n.Position(), "failed to find struct address for type '%s'", newType.Name)
		}
		return

	case *ast.BuiltinType:
		// Add the type directly
		newType.Struct = processBasic(v.Type)
		if newType.Struct == nil {
			context.AddTypeError(n.Position(), "failed to find basic type for type '%s'", newType.Name)
		}
		return

	case *ast.Typedef:
		// Don't need to descend further. We care about the structs used, not the typedefs.
		return

	case *ast.Record:
		// A record is a struct def - this is what we need
		address = &v.Addr
		return

	default:
		// Don't descend any further
		return
	}

}

// ProcessTypes parses the AST information and converts it into lists of CStructs and CTypes. If errors are encountered
// during processing, they are recorded in the context and processing continues.
func ProcessTypes(context *errors.ParseContext, recordDeclarations []*ast.RecordDecl,
	typeDeclarations []*ast.TypedefDecl) (structs []*CStruct, types []*CType) {

	structsByAddress := make(map[ast.Address]*CStruct)
	structsByName := make(map[string]*CStruct)
	typesByName := make(map[string]*CType)

	// Record declarations are structs. During the first pass, we record each struct we find, ignoring fields.
	for _, r := range recordDeclarations {
		if r.Kind != "struct" {
			continue
		}

		newStruct := &CStruct{Name: r.Name}
		structsByAddress[r.Address()] = newStruct

		if newStruct.Name != "" {
			structsByName[newStruct.Name] = newStruct
		}

		structs = append(structs, newStruct)
	}

	// Now process typedefs, linking back to the underlying struct (counting pointers as we go).
	for _, t := range typeDeclarations {
		if strings.HasPrefix(t.Name, "__") {
			// skip implementation types
			continue
		}

		newType := &CType{Name: t.Name}

		if len(t.ChildNodes) < 1 {
			context.AddTypeError(t.Pos, "Cannot find child nodes of type %s", t.Name)
			continue
		}

		mark := context.Mark()
		structAddr := processTypeNode(context, t.ChildNodes[0], newType)
		if context.HasErrors(mark) {
			continue
		}

		// We may not find a struct for certain types (so far, it seems to be some magic built-in types that
		// trip us up. So print a warning in that case and continue
		if structAddr == nil && newType.Struct == nil {
			context.AddTypeError(t.Pos, "no struct found for type %s", t.Name)
			continue
		}

		// If newType is a typedef for a basic type (e.g. char, int), then newType.Struct will be non-nill
		// and we don't need to look for the struct by its address.
		if newType.Struct == nil {
			s, found := structsByAddress[*structAddr]
			if !found {
				context.AddTypeError(t.Pos, "cannot find struct from AST address, for type %s", t.Name)
				continue
			}

			newType.Struct = s

			// Update struct type name, if needed
			if s.TypeName == "" && newType.PointerCount == 0 {
				s.TypeName = newType.Name
			}
		}

		types = append(types, newType)
		typesByName[t.Name] = newType
	}

	// Now we know about every available struct and type, re-process the struct declarations to add fields
	for _, r := range recordDeclarations {
		if r.Kind != "struct" {
			// We have already printed a warning message about this
			continue
		}

		// Find the struct again
		s := structsByAddress[r.Address()]

		for _, n := range r.ChildNodes {
			f, ok := n.(*ast.FieldDecl)
			if !ok {
				continue
			}

			// Is this an array type?
			nonArrayType, arrayCount := getArrayLength(f.Type)

			// Count any pointers
			nonPointerType, numPointers := stripPointers(nonArrayType)

			// Check if this is a basic type
			childStruct := processBasic(nonPointerType)

			// If not, look in the known types
			if childStruct == nil {
				if strings.HasPrefix(nonPointerType, "struct ") {
					childStruct = structsByName[strings.TrimPrefix(nonPointerType, "struct ")]
				} else {
					childType := typesByName[nonPointerType]
					if childType != nil {
						childStruct = childType.Struct
						numPointers += childType.PointerCount
					}
				}

				if childStruct == nil {
					context.AddTypeError(f.Pos, "cannot find type '%s' referenced in field '%s' of struct '%s'",
						f.Type, f.Name, s.Name)
					continue
				}
			}

			s.Fields = append(s.Fields, &CStructInstance{
				Name:         f.Name,
				Struct:       childStruct,
				PointerCount: numPointers,
				ArrayCount:   arrayCount,
			})
		}
	}

	return
}

var charStruct = &CStruct{Name: "char", Basic: true}
var signedCharStruct = &CStruct{Name: "signed char", Basic: true}
var unsignedCharStruct = &CStruct{Name: "unsigned char", Basic: true}
var shortStruct = &CStruct{Name: "short", Basic: true}
var unsignedShortStruct = &CStruct{Name: "unsigned short", Basic: true}
var intStruct = &CStruct{Name: "int", Basic: true}
var unsignedStruct = &CStruct{Name: "unsigned", Basic: true}
var longStruct = &CStruct{Name: "long", Basic: true}
var unsignedLongStruct = &CStruct{Name: "unsigned long", Basic: true}
var longLongStruct = &CStruct{Name: "long long", Basic: true}
var unsignedLongLongStruct = &CStruct{Name: "unsigned long long", Basic: true}
var floatStruct = &CStruct{Name: "float", Basic: true}
var doubleStruct = &CStruct{Name: "double", Basic: true}
var voidStruct = &CStruct{Name: "void", Basic: true}

func processBasic(ctype string) *CStruct {

	switch ctype {
	case "char":
		return charStruct

	case "signed char":
		return signedCharStruct

	case "unsigned char":
		return unsignedCharStruct

	case "short", "short int", "signed short", "signed short int":
		return shortStruct

	case "unsigned short", "unsigned short int":
		return unsignedShortStruct

	case "int", "signed", "signed int":
		return intStruct

	case "unsigned", "unsigned int":
		return unsignedStruct

	case "long", "long int", "signed long", "signed long int":
		return longStruct

	case "unsigned long", "unsigned long int":
		return unsignedLongStruct

	case "long long", "long long int", "signed long long", "signed long long int":
		return longLongStruct

	case "unsigned long long", "unsigned long long int":
		return unsignedLongLongStruct

	case "float":
		return floatStruct

	case "double":
		return doubleStruct

	case "void":
		return voidStruct

	default:
		return nil
	}
}

func stripPointers(typeName string) (nonPointerType string, pointerCount uint) {
	pointerCount = uint(strings.Count(typeName, "*"))
	nonPointerType = strings.TrimSpace(strings.Replace(typeName, "*", "", -1))
	return
}

var arrayRegex = regexp.MustCompile(`(.*) \[(\d+)]`)

// getArrayLength parses a type looking for a "foo [N]" pattern. If found, it returns
// "foo" and N respectively.
func getArrayLength(typeName string) (nonArrayType string, arrayCount uint) {
	if arrayRegex.MatchString(typeName) {
		subMatches := arrayRegex.FindStringSubmatch(typeName)

		arrayCountInt, err := strconv.Atoi(subMatches[2])
		if err != nil {
			log.Panicf("%s", err) // Can't fail, due to regex match
		}
		return subMatches[1], uint(arrayCountInt)
	}

	nonArrayType = typeName
	return
}
