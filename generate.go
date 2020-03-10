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
	"os"
	"path/filepath"
)

// Apply executes all macros within the current context.
func Apply() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	dir = findRoot(dir)
	fmt.Println(dir)
	files, err := parse(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		f.parseComments()
	}
	for _, f := range files {
		fmt.Println(f.String())
	}
	return nil
}

func findRoot(dir string) string {
	for len(dir) > 0 {
		gomod := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(gomod); err != nil {
			dir = filepath.Dir(dir)
		}
		return dir
	}
	return dir
}

func MustApply() {
	err := Apply()
	if err != nil {
		panic(err)
	}
}
