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
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thales-e-security/header2go/translate/errors"
)

const testDir = "testdata"

// Tests are defined individually to aid debugging.

func TestBasicTypes(t *testing.T) {
	runTest(t, "01_basic_types")
}

func TestSimpleStructs(t *testing.T) {
	runTest(t, "02_simple_structs")
}

func TestSimpleStructsWithReturnTypes(t *testing.T) {
	runTest(t, "03_simple_structs_with_return_types")
}

func TestSimplePointerArgs(t *testing.T) {
	runTest(t, "04_simple_pointer_args")
}

func TestTypeDefsOfBasicTypes(t *testing.T) {
	runTest(t, "05_typedefs_of_basic_types")
}

func TestStructsWithKeywords(t *testing.T) {
	runTest(t, "06_structs_with_keywords")
}

func TestParamsWithKeywords(t *testing.T) {
	runTest(t, "07_params_with_keywords")
}

func TestPointersToStructs(t *testing.T) {
	runTest(t, "08_pointers_to_structs_with_basic_types")
}

func TestFixedLengthArrays(t *testing.T) {
	runTest(t, "09_fixed_length_arrays")
}

func runTest(t *testing.T, testName string) {
	/*
		For each directory <X>, we expect to find:

			<X>/defs.go
			<X>/input.h
			<X>/main.go
	*/

	t.Logf("Running sub-test: %s", testName)
	outDir, err := ioutil.TempDir("", "header2go")
	defer os.RemoveAll(outDir)

	context := errors.NewParseContext()
	mark := context.Mark()
	err = Process(context, path.Join(testDir, testName, "input.h"), outDir)
	require.NoError(t, err)
	require.False(t, context.HasErrors(mark))

	// Compare main.go files (assumes both are gofmt-ed)
	expected, err := ioutil.ReadFile(path.Join(testDir, testName, "main.go"))
	require.NoError(t, err)

	actual, err := ioutil.ReadFile(path.Join(outDir, "main.go"))
	require.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))

	// Compare defs.go files (assumes both are gofmt-ed)
	expected, err = ioutil.ReadFile(path.Join(testDir, testName, "defs.go"))
	require.NoError(t, err)

	actual, err = ioutil.ReadFile(path.Join(outDir, "defs.go"))
	require.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}
