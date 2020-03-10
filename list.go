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
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// parse returns all files recursively
func parse(dir string) ([]*File, error) {
	fset := token.NewFileSet()
	files, err := listFiles(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to list files in %s: %w", dir, err)
	}
	for _, f := range files {
		f.fset = fset
		node, err := parser.ParseFile(fset, f.fileName, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("unable to parse %s: %w", f.fileName, err)
		}
		f.ast = node
	}
	return files, nil
}

// listFiles returns all *.go files recursively, as long as they are not hidden or in hidden folders
func listFiles(dir string) ([]*File, error) {
	var res []*File
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Mode().IsRegular() && strings.HasSuffix(strings.ToLower(info.Name()), ".go") {
			res = append(res, &File{fileName: path, info: info})
		}
		return nil
	})
	return res, err
}
