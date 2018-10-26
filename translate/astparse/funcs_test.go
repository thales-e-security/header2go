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

func TestParseFunc(t *testing.T) {
	testFile, err := filepath.Abs("testdata/types1.h")
	require.NoError(t, err)

	functions, records, typedefs, err := ReadAst(testFile)
	require.NoError(t, err)

	context := errors.NewParseContext()

	mark := context.Mark()
	structs, types := ProcessTypes(context, records, typedefs)
	require.False(t, context.HasErrors(mark))

	mark = context.Mark()
	funcs, usedTypes := ProcessFuncs(context, functions, structs, types)
	require.False(t, context.HasErrors(mark))

	require.Len(t, usedTypes, 3)
	require.Len(t, funcs, 2)

	// Types are returned in the order they are defined, not in the order of function parameters.
	assert.Equal(t, "Struct1", usedTypes[0].Name)
	assert.Equal(t, "", usedTypes[1].Name)
	assert.Equal(t, "Struct3", usedTypes[2].Name)

	// Function A
	assert.Equal(t, "functionA", funcs[0].Name)
	assert.Nil(t, funcs[0].ReturnType)

	assert.Equal(t, "Struct3", funcs[0].Params[0].Struct.Name)
	assert.Equal(t, "a1", funcs[0].Params[0].Name)
	assert.Equal(t, uint(1), funcs[0].Params[0].PointerCount)

	assert.Equal(t, "long", funcs[0].Params[1].Struct.Name)
	assert.True(t, funcs[0].Params[1].Struct.Basic)
	assert.Equal(t, "a2", funcs[0].Params[1].Name)
	assert.Equal(t, uint(0), funcs[0].Params[1].PointerCount)

	// Function B
	assert.Equal(t, "functionB", funcs[1].Name)
	assert.Equal(t, "char", funcs[1].ReturnType.Struct.Name)
	assert.Equal(t, uint(1), funcs[1].ReturnType.PointerCount)

	assert.Equal(t, "Struct3", funcs[1].Params[0].Struct.Name)
	assert.Equal(t, "a1", funcs[1].Params[0].Name)
	assert.Equal(t, uint(0), funcs[1].Params[0].PointerCount)

	assert.Equal(t, "", funcs[1].Params[1].Struct.Name)
	assert.Equal(t, "a2", funcs[1].Params[1].Name)
	assert.Equal(t, uint(1), funcs[1].Params[1].PointerCount)
}
