// Copyright © 2018 Christoph Petrausch
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exploration

import (
	"os"
	"path/filepath"
)

// InitialFiles return all files and directories in the tree below rootDir and the rootDir itself
func InitialFiles(rootDir string) (directories []string, files []string, err error) {
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil && info == nil {
			return nil
		}
		if info.IsDir() {
			directories = append(directories, path)
		} else {
			files = append(files, path)
		}
		return nil
	}
	err = filepath.Walk(rootDir, walkFunc)
	return directories, files, err
}
