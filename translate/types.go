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
	"fmt"
	"strings"

	"github.com/thales-e-security/header2go/translate/astparse"
)

// keywords is a map of all the language keywords
var keywords = map[string]bool{
	"break":       true,
	"default":     true,
	"func":        true,
	"interface":   true,
	"select":      true,
	"case":        true,
	"defer":       true,
	"go":          true,
	"map":         true,
	"struct":      true,
	"chan":        true,
	"else":        true,
	"goto":        true,
	"package":     true,
	"switch":      true,
	"const":       true,
	"fallthrough": true,
	"if":          true,
	"range":       true,
	"type":        true,
	"continue":    true,
	"for":         true,
	"import":      true,
	"return":      true,
	"var":         true,
}

// templateType decorates the AST struct type with methods for generating template outputs
type templateType astparse.CStruct

// CanonicalName is equal to TypeName, if not empty, otherwise StructName
func (t *templateType) CanonicalName() string {
	if t.TypeName != "" {
		return t.TypeName
	}
	return t.Name
}

// GoStructName is a suitable name for a Go type representing this C type.
func (t *templateType) GoStructName() string {
	if t.Basic {
		return basicToGo(t.Name)
	}

	if t.TypeName != "" {
		return capitalise(t.TypeName)
	}
	return capitalise(t.Name)
}

// GoArrayTypeName returns the type name to use when we have an array of this type.
func (t *templateType) GoArrayTypeName() string {
	return capitalise(t.GoStructName()) + "Array"
}

// CgoStructName is a name safe for using with cgo (i.e. C.foo format). It uses the TypeName, if present,
// otherwise the StructName prefixed with "struct_". In either case, the "C." is prefixed.
func (t *templateType) CgoStructName() string {
	return "C." + strings.Replace(t.CanonicalName(), " ", "_", -1)
}

// ConversionFunction is a suitable function name for converting this C type to a Go type.
func (t *templateType) ConversionFunction() string {
	if t.Basic {
		return t.GoStructName()
	}

	return "convertTo" + t.GoStructName()
}

// RConversionFunction is a suitable function name for converting the Go type to this C type to a Go type.
func (t *templateType) RConversionFunction() string {
	if t.Basic {
		return t.CgoStructName()
	}

	return "convertFrom" + t.GoStructName()
}

// EmptyValue is a suitable nil value for this type (in Go). For structs it is just "<struct>{}", while for
// basic types it is the normal zero/nil value.
func (t *templateType) EmptyValue() string {
	if t.Basic {
		return "0"
	}

	return t.GoStructName() + "{}"
}

// GetFields returns the fields of this struct as templateTypeInstance values. (Use this rather than
// accessing the Fields field).
func (t *templateType) GetFields() []*templateTypeInstance {
	var res []*templateTypeInstance

	for _, f := range t.Fields {
		tti := templateTypeInstance(*f)
		res = append(res, &tti)
	}
	return res
}

// templateTypeInstance decorates a CStructInstance with information needed for templates.
type templateTypeInstance astparse.CStructInstance

// GoName contains a go-friendly variation of Name.
func (t templateTypeInstance) GoName() string {
	return capitalise(t.Name)
}

func (t templateTypeInstance) Type() *templateType {
	tt := templateType(*t.Struct)
	return &tt
}

// Pointers returns a string of "*" values to turn pointer types into pointers.
func (t templateTypeInstance) Pointers() string {
	return strings.Repeat("*", int(t.PointerCount))
}

// ParamAsArg returns the string to be used when the parameter is passed to the Go function.
// This function is used to avoid some more complicated logic in the template.
func (t templateTypeInstance) ParamAsArg() string {
	if t.PointerCount > 0 {
		// We will have created a variable of the correct type already
		return t.Name + "Array"
	}

	// Otherwise we need a normal conversion
	return fmt.Sprintf("%s(%s)", t.Type().ConversionFunction(), t.SafeName())
}

// GoArgType returns the type to use in the Go version of the function.
func (t templateTypeInstance) GoArgType() string {
	if t.PointerCount > 0 {
		return "*" + t.Type().GoArrayTypeName()
	}
	return t.Type().GoStructName()
}

func (t templateTypeInstance) SafeName() string {

	if keywords[t.Name] {
		return "_" + t.Name
	}

	return t.Name
}

// ArrayAwareConversion returns a C->Go conversion function that takes into account whether this type
// has been used in an array.
func (t templateTypeInstance) ArrayAwareConversion() string {
	if t.ArrayCount == 0 {
		return t.Type().ConversionFunction()
	}

	return fmt.Sprintf("convertTo%sArray%d", t.Type().GoStructName(), t.ArrayCount)
}

// templateFunc describes a function found in the C header file.
type templateFunc astparse.CFuncDef

// GoName is a go-friendly variant of Name.
func (t *templateFunc) GoName() string {
	return "go" + capitalise(t.Name)
}

// GetParams returns the parameters as templateTypeInstance values. Use this instead of
// direct access to the Params field.
func (t *templateFunc) GetParams() []*templateTypeInstance {
	var res []*templateTypeInstance

	for _, p := range t.Params {
		tti := templateTypeInstance(*p)
		res = append(res, &tti)
	}
	return res
}

// GetReturnType returns the function return type, or nil for void functions.
func (t *templateFunc) GetReturnType() *templateTypeInstance {
	if t.ReturnType == nil {
		return nil
	}

	tti := templateTypeInstance(*t.ReturnType)
	return &tti
}

type templateArrayType struct {
	Type       *templateType
	ArrayCount uint
}
