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

package errors

import "github.com/elliotchance/c2go/ast"

// A ParseContext contains error information about a parsing activity.
type ParseContext struct {
	typeErrors     []SourceError
	functionErrors []SourceError
	count          int
	marks          []int
}

// NewParseContext creates a new context.
func NewParseContext() *ParseContext {
	return &ParseContext{}
}

// AddTypeError adds a new error relating to type parsing.
func (c *ParseContext) AddTypeError(pos ast.Position, format string, args ...interface{}) {
	c.typeErrors = append(c.typeErrors, NewSourceError(pos, format, args...))
	c.count++
}

// AddFunctionError adds a new error relating to function parsing.
func (c *ParseContext) AddFunctionError(pos ast.Position, format string, args ...interface{}) {
	c.functionErrors = append(c.functionErrors, NewSourceError(pos, format, args...))
	c.count++
}

// Mark provides an opaque reference to the current context. Later, this reference can be used
// in HasErrors to determine if the context has changed.
func (c *ParseContext) Mark() (mark int) {
	mark = len(c.marks)
	c.marks = append(c.marks, c.count)
	return
}

// HasErrors returns true if an error has been recorded since the Mark was called.
func (c *ParseContext) HasErrors(mark int) bool {
	return c.count > c.marks[mark]
}

func cloneErrors(src []SourceError) []SourceError {
	res := make([]SourceError, len(src))
	copy(res, src)
	return res
}

// TypeErrors returns the type errors (possibly empty).
func (c *ParseContext) TypeErrors() []SourceError {
	return cloneErrors(c.typeErrors)
}

// FunctionErrors returns the function errors (possibly empty).
func (c *ParseContext) FunctionErrors() []SourceError {
	return cloneErrors(c.functionErrors)
}
