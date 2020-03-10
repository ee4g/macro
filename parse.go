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
	"strings"
)

const macroImportPrefix = "#[import]"

const stateFindFuncBodyStart = 0
const stateFindFuncBodyEnd = 1

func commentText(g *ast.CommentGroup) string {
	var lines []string
	for _, c := range g.List {
		c := c.Text
		switch c[1] {
		case '/':
			c = "  " + c[2:] + "\n" // need to add white space so that column count matches correctly
		case '*':
			c = "  " + c[2:len(c)-2] // need to add white space so that column count matches correctly
		}
		lines = append(lines, c)
	}
	return strings.Join(lines, "")
}

func (f *File) parseComments() {
	// this is already a workaround because the go ast does not contain floating comments, which is exactly what
	// we require, because we want to be excluded from attached comments
	for _, cgroup := range f.ast.Comments {
		state := stateFindFuncBodyStart
		bracketBalance := 0
		macro := &Block{Parent: f}
		cleanSplitLines := strings.Split(commentText(cgroup), "\n")

		for lineNo, line := range cleanSplitLines {
			tmpLine := &strings.Builder{}
			for i := 0; i < len(line); i++ {
				switch state {
				case stateFindFuncBodyStart:
					if line[i] == '#' && i < len(line)-1 && line[i+1] == '[' {
						state = stateFindFuncBodyEnd
						i++
					}
				case stateFindFuncBodyEnd:
					if line[i] == '[' {
						tmpLine.WriteByte(line[i])
						bracketBalance++
					} else
					if line[i] == ']' {
						if bracketBalance == 0 {
							state = stateFindFuncBodyStart
						} else {
							tmpLine.WriteByte(line[i])
							bracketBalance--
						}
					} else {
						tmpLine.WriteByte(line[i])
					}
				default:
					panic("invalid state")
				}

			}

			if tmpLine.Len() > 0 {
				stmt := strings.TrimSpace(tmpLine.String())

				macro.Lines = append(macro.Lines, &Line{
					Pos: Position{
						File:   f,
						Line:   lineNo + f.fset.Position(cgroup.Pos()).Line,
						Column: strings.Index(cleanSplitLines[lineNo], stmt) + f.fset.Position(cgroup.Pos()).Column,
					},
					SrcText:   cleanSplitLines[lineNo],
					Statement: stmt,
				})
			}

		}

		if len(macro.Lines) > 0 {
			importMacro := false
			for _, macroLine := range macro.Lines {
				switch macroLine.Statement {
				case "import":
					importMacro = true

					f.imports = append(f.imports, &Import{
						Pos:       macroLine.Pos,
						Statement: strings.TrimSpace(macroLine.SrcText[strings.Index(macroLine.SrcText, macroImportPrefix)+len(macroImportPrefix):]),
						SrcText:   macroLine.SrcText,
					})
				}

			}
			if !importMacro {
				macro.Annotated = f.nextTypeOrFuncNode(f.ast, cgroup.Pos())
				f.macros = append(f.macros, macro)

			}
		}
		cgroup.Pos()
	}
}

func (f *File) nextTypeOrFuncNode(root ast.Node, offset token.Pos) ast.Node {
	var res ast.Node
	searchPos := f.fset.Position(offset)
	ast.Inspect(root, func(node ast.Node) bool {
		if node == nil {
			return true
		}
		if res != nil{
			return false
		}
		npos := f.fset.Position(node.Pos())
		if npos.Filename == searchPos.Filename {
			if npos.Line > searchPos.Line {
				switch t := node.(type) {
				case *ast.GenDecl:
					if t.Tok == token.TYPE {
						res = t
						return false
					}
				case *ast.FuncDecl:
					res = t
					return false
				}
				return false
			}
		}

		return true
	})
	return res
}
