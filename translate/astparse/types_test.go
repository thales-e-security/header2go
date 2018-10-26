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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thales-e-security/header2go/translate/errors"
)

// TestTypeParsing parses a carefully constructed header file that should trigger most of the interesting
// situations.
func TestTypeParsing(t *testing.T) {
	testFile, err := filepath.Abs("testdata/types1.h")
	require.NoError(t, err)

	_, records, typedefs, err := ReadAst(testFile)
	require.NoError(t, err)

	context := errors.NewParseContext()
	mark := context.Mark()
	structs, types := ProcessTypes(context, records, typedefs)
	require.False(t, context.HasErrors(mark))

	// Validate lengths
	require.Len(t, structs, 3)

	for _, v := range types {
		t.Log(v.Name, v.Struct.Name)
	}

	require.Len(t, types, 6)

	// Validate structs
	structsByName := make(map[string]*CStruct)
	for _, s := range structs {
		structsByName[s.Name] = s
	}

	s1 := structsByName["Struct1"]
	assert.False(t, s1.Basic)
	assert.Equal(t, "a", s1.Fields[0].Name)
	assert.Equal(t, "int", s1.Fields[0].Struct.Name)
	assert.True(t, s1.Fields[0].Struct.Basic)
	assert.Nil(t, s1.Fields[0].Struct.Fields)

	assert.Equal(t, "b", s1.Fields[1].Name)
	assert.Equal(t, "long", s1.Fields[1].Struct.Name)
	assert.True(t, s1.Fields[1].Struct.Basic)
	assert.Nil(t, s1.Fields[1].Struct.Fields)

	s2 := structsByName[""]
	assert.False(t, s2.Basic)
	assert.Equal(t, "a", s2.Fields[0].Name)
	assert.Equal(t, s1, s2.Fields[0].Struct)
	assert.Equal(t, "b", s2.Fields[1].Name)
	assert.Equal(t, "int", s2.Fields[1].Struct.Name)
	assert.Equal(t, uint(1), s2.Fields[1].PointerCount)

	s3 := structsByName["Struct3"]
	assert.False(t, s3.Basic)
	assert.Equal(t, "a", s3.Fields[0].Name)
	assert.Equal(t, uint(1), s3.Fields[0].PointerCount)
	assert.Equal(t, s1, s3.Fields[0].Struct)

	// Process types
	typesByName := make(map[string]*CType)
	for _, t := range types {
		typesByName[t.Name] = t
	}

	t1 := typesByName["Type2"]
	assert.Equal(t, uint(0), t1.PointerCount)
	assert.Equal(t, s2, t1.Struct)

	t2 := typesByName["Type1"]
	assert.Equal(t, uint(0), t2.PointerCount)
	assert.Equal(t, s1, t2.Struct)

	t3 := typesByName["PtrType1"]
	assert.Equal(t, uint(1), t3.PointerCount)
	assert.Equal(t, s1, t3.Struct)

	t4 := typesByName["Type3"]
	assert.Equal(t, uint(0), t4.PointerCount)
	assert.Equal(t, s3, t4.Struct)

	t5 := typesByName["PtrType3"]
	assert.Equal(t, uint(1), t5.PointerCount)
	assert.Equal(t, s3, t5.Struct)

	t6 := typesByName["AnotherPtrType3"]
	assert.Equal(t, uint(1), t6.PointerCount)
	assert.Equal(t, s3, t6.Struct)
}

func TestDoubleTypeDef(t *testing.T) {
	testFile, err := filepath.Abs("testdata/types2.h")
	require.NoError(t, err)

	_, records, typedefs, err := ReadAst(testFile)
	require.NoError(t, err)

	context := errors.NewParseContext()
	mark := context.Mark()
	ProcessTypes(context, records, typedefs)
	require.False(t, context.HasErrors(mark))
}

func TestArrayParsing(t *testing.T) {
	testFile, err := filepath.Abs("testdata/fixed_arrays.h")
	require.NoError(t, err)

	_, records, typedefs, err := ReadAst(testFile)
	require.NoError(t, err)

	context := errors.NewParseContext()
	mark := context.Mark()
	structs, _ := ProcessTypes(context, records, typedefs)
	require.False(t, context.HasErrors(mark))

	// Validate lengths
	require.Len(t, structs, 1)

	assert.Equal(t, uint(0), structs[0].Fields[0].ArrayCount)
	assert.Equal(t, uint(15), structs[0].Fields[1].ArrayCount)
}
