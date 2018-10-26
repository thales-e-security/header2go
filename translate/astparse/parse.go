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

// Portions of this file have been copied from https://github.com/elliotchance/c2go and are copyrighted by
// Elliot Chance, Izyumov Konstantin and Christophe Kamphaus under an MIT license.

package astparse

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/preprocessor"
	"github.com/elliotchance/c2go/program"
	"github.com/pkg/errors"
)

// ReadAst processes the header file and converts it into the AST types used by the c2go
// project.
func ReadAst(headerFile string) ([]*ast.FunctionDecl, []*ast.RecordDecl,
	[]*ast.TypedefDecl, error) {

	// Code taken from github.com/elliotchance/c2go/main.go, with some modifications.

	// Preprocess code and save to temp file
	pp, comments, includes, err := preprocessor.Analyze([]string{headerFile}, []string{})
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to preprocess '%s'", headerFile)
	}

	dir, err := ioutil.TempDir("", "header2go")
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to open temp dir")
	}
	defer os.RemoveAll(dir)

	ppFilePath := path.Join(dir, "pp.c")
	err = ioutil.WriteFile(ppFilePath, pp, 0644)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to write to file '%s'", ppFilePath)
	}

	// Generate JSON from AST
	// #nosec G204 - ppFilePath is not based on user input.
	astPP, err := exec.Command("clang", "-Xclang", "-ast-dump",
		"-fsyntax-only", "-fno-color-diagnostics", ppFilePath).Output()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to generate JSON from AST")
	}

	p := program.NewProgram()
	p.Comments = comments
	p.IncludeHeaders = includes

	// Converting to nodes
	nodes := convertLinesToNodes(strings.Split(string(astPP), "\n"))

	// build tree
	tree := buildTree(nodes, 0)
	ast.FixPositions(tree)
	p.SetNodes(tree)

	// skip call to ast.RepairCharacterLiteralsFromSource

	// Find functions in tree
	var f func(ast.Node)

	var functions []*ast.FunctionDecl
	var recordDeclarations []*ast.RecordDecl
	var typeDeclarations []*ast.TypedefDecl

	f = func(node ast.Node) {
		switch n := node.(type) {
		case *ast.FunctionDecl:
			functions = append(functions, n)
		case *ast.RecordDecl:
			recordDeclarations = append(recordDeclarations, n)
		case *ast.TypedefDecl:
			typeDeclarations = append(typeDeclarations, n)
		}

		for _, n := range node.Children() {
			f(n)
		}
	}

	for _, n := range tree {
		f(n)
	}

	return functions, recordDeclarations, typeDeclarations, nil
}

type treeNode struct {
	indent int
	node   ast.Node
}

// Function taken from github.com/elliotchance/c2go/main.go
func convertLinesToNodes(lines []string) []treeNode {
	nodes := make([]treeNode, len(lines))
	var counter int
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// It is tempting to discard null AST nodes, but these may
		// have semantic importance: for example, they represent omitted
		// for-loop conditions, as in for(;;).
		line = strings.Replace(line, "<<<NULL>>>", "NullStmt", 1)
		trimmed := strings.TrimLeft(line, "|\\- `")
		node := ast.Parse(trimmed)
		indentLevel := (len(line) - len(trimmed)) / 2
		nodes[counter] = treeNode{indentLevel, node}
		counter++
	}
	nodes = nodes[0:counter]

	return nodes
}

// Function taken from github.com/elliotchance/c2go/main.go
func buildTree(nodes []treeNode, depth int) []ast.Node {
	if len(nodes) == 0 {
		return []ast.Node{}
	}

	// Split the list into sections, treat each section as a tree with its own
	// root.
	sections := [][]treeNode{}
	for _, node := range nodes {
		if node.indent == depth {
			sections = append(sections, []treeNode{node})
		} else {
			sections[len(sections)-1] = append(sections[len(sections)-1], node)
		}
	}

	results := []ast.Node{}
	for _, section := range sections {
		slice := []treeNode{}
		for _, n := range section {
			if n.indent > depth {
				slice = append(slice, n)
			}
		}

		children := buildTree(slice, depth+1)
		for _, child := range children {
			section[0].node.AddChild(child)
		}
		results = append(results, section[0].node)
	}

	return results
}
