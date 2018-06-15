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
	"testing"

	"path"

	"github.com/stretchr/testify/assert"
)

func TestInitialFiles(t *testing.T) {
	tests := []struct {
		name                string
		dir                 string
		cleanDir            bool
		filesToTouch        []touchFile
		expectedFiles       []string
		expectedDirectories []string
		expectedError       error
	}{
		{
			name:                "empty dir",
			dir:                 createTempDir(t),
			cleanDir:            true,
			expectedDirectories: []string{""},
		},
		{
			name:     "not existing dir",
			dir:      "/tmp/foo-bar",
			cleanDir: false,
		},
		{
			name:     "dir with file and subdir",
			dir:      createTempDir(t),
			cleanDir: true,
			filesToTouch: []touchFile{
				{
					name:  "foo",
					isDir: true,
				},
				{
					name:  "foo/bar",
					isDir: false,
				},
				{
					name:  "baz",
					isDir: false,
				},
			},
			expectedDirectories: []string{"", "foo"},
			expectedFiles:       []string{"baz", "foo/bar"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.cleanDir {
				defer os.RemoveAll(test.dir)
			}
			touchFiles(t, test.dir, test.filesToTouch)
			dirs, files, err := InitialFiles(test.dir)
			joinPathsWithTempFile(test.dir, test.expectedFiles)
			joinPathsWithTempFile(test.dir, test.expectedDirectories)
			assert.Equal(t, test.expectedFiles, files)
			assert.Equal(t, test.expectedDirectories, dirs)
			assert.Equal(t, test.expectedError, err)

		})
	}
}

func joinPathsWithTempFile(testDir string, paths []string) {
	for i := range paths {
		paths[i] = path.Join(testDir, paths[i])
	}
}
