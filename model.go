/*
 * Copyright 2020 Torben Schinke
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package macro

import (
	"go/ast"
	"go/token"
	"os"
	"strconv"
	"strings"
)

// A Position is the absolute line number starting at 1 and column also starting at 1, just as IDEs and
// the compiler does.
type Position struct {
	File   *File
	Line   int
	Column int
}

// A File contains all imports and macros together
type File struct {
	fset     *token.FileSet
	info     os.FileInfo // info about the source file
	fileName string      // the original file name
	ast      *ast.File   // the parsed AST
	imports  []*Import   // the declared and parsed imports
	macros   []*Block    // the declared and parsed macro blocks
}

func (f *File) String() string {
	sb := &strings.Builder{}
	sb.WriteString(f.fileName)
	sb.WriteString(":\n")
	sb.WriteString("imports:\n")
	for _, i := range f.imports {
		sb.WriteString(i.String())
		sb.WriteString("\n")
	}
	sb.WriteString("macros:\n")
	for _, i := range f.macros {
		sb.WriteString(i.String())
		sb.WriteString("\n")
	}
	return sb.String()
}

// Import for packages
type Import struct {
	Pos       Position
	Statement string
	SrcText   string
}

func (i Import) String() string {
	return i.Pos.File.info.Name() + ":" + strconv.Itoa(i.Pos.Line) + ":" + strconv.Itoa(i.Pos.Column) + ":" + i.Statement
}

// Block is a collection of lines
type Block struct {
	Lines     []*Line
	Annotated ast.Node
	Parent    *File
}

func (b *Block) String() string {
	sb := &strings.Builder{}

	annotatedPos := b.Parent.fset.Position(b.Annotated.Pos())
	annotatedLoc := annotatedPos.Filename + ":" + strconv.Itoa(annotatedPos.Line) + ":" + strconv.Itoa(annotatedPos.Column)

	switch t := b.Annotated.(type) {
	case *ast.FuncDecl:
		sb.WriteString("// macro for function '" + t.Name.Name + "' defined at " + annotatedLoc + "\n")
	case *ast.GenDecl:
		typeSpec := t.Specs[0].(*ast.TypeSpec)
		sb.WriteString("// macro for type '" + typeSpec.Name.Name + "' defined at " + annotatedLoc + "\n")
	}
	sb.WriteString("func macro(){\n")
	for _, line := range b.Lines {
		// //line filename:line:col
		locPrint := line.Pos.File.info.Name() + ":" + strconv.Itoa(line.Pos.Line) + ":" + strconv.Itoa(line.Pos.Column)
		sb.WriteString("  " + line.Statement + " // defined at " + locPrint)
		sb.WriteString("\n")
	}
	sb.WriteString("}")
	return sb.String()
}

// Line
type Line struct {
	Pos       Position
	Statement string
	SrcText   string
}
