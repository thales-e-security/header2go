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

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoodConfig(t *testing.T) {
	cfg, err := FromFile("testdata/good-config.toml")
	require.NoError(t, err)

	expected := ParseConfig{
		VoidParam: []VoidParam{
			{Function: "doSomething", Parameter: "val1", ReplaceWith: "SomeType"},
		},
		VoidField: []VoidField{
			{TypeName: "struct someStruct", Field: "someField", ReplaceWith: "struct someStruct"},
		},
	}

	assert.Equal(t, expected, cfg)
}

func TestVoidParamFunctionMissing(t *testing.T) {
	input := `
[[voidParam]]
#function = "doSomething"
parameter = "val1"
replaceWith = "SomeType" 
`
	_, err := FromString(input)
	require.Error(t, err)
}

func TestVoidParamParameterMissing(t *testing.T) {
	input := `
[[voidParam]]
function = "doSomething"
#parameter = "val1"
replaceWith = "SomeType" 
`
	_, err := FromString(input)
	require.Error(t, err)
}

func TestVoidParamReplaceWithMissing(t *testing.T) {
	input := `
[[voidParam]]
function = "doSomething"
parameter = "val1"
#replaceWith = "SomeType" 
`
	_, err := FromString(input)
	require.Error(t, err)
}

func TestVoidFieldTypeNameMissing(t *testing.T) {
	input := `
[[voidField]]
#typeName = "struct someStruct"
field = "someField"
replaceWith = "struct someStruct" 
`
	_, err := FromString(input)
	require.Error(t, err)
}

func TestVoidFieldFieldMissing(t *testing.T) {
	input := `
[[voidField]]
typeName = "struct someStruct"
#field = "someField"
replaceWith = "struct someStruct" 
`
	_, err := FromString(input)
	require.Error(t, err)
}

func TestVoidFieldReplaceWithMissing(t *testing.T) {
	input := `
[[voidField]]
typeName = "struct someStruct"
field = "someField"
#replaceWith = "struct someStruct" 
`
	_, err := FromString(input)
	require.Error(t, err)
}

func TestNonsenseConfig(t *testing.T) {
	_, err := FromString("not valid")
	require.Error(t, err)
}

func TestBadFilePath(t *testing.T) {
	_, err := FromFile("this does not exist")
	require.Error(t, err)
}
