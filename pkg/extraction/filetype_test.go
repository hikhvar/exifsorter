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

package extraction

import (
	"fmt"
	"path"
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestIsVideoOrImage(t *testing.T) {
	tests := []struct {
		name          string
		fileOrVideo   bool
		expectedError string
	}{
		{
			name:        "sample1.JPG",
			fileOrVideo: true,
		},
		{
			name:        "sample2.mp4",
			fileOrVideo: true,
		},
		{
			name:        "sample3.txt",
			fileOrVideo: false,
		},
		{
			name:          "sample-not-exist",
			fileOrVideo:   false,
			expectedError: "could not open file to determine file type: open %s: no such file or directory",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fileUnderTest := fixturePath(test.name)
			is, err := IsVideoOrImage(fileUnderTest)
			assert.Equal(t, test.fileOrVideo, is)
			if test.expectedError != "" {
				assert.EqualError(t, err, fmt.Sprintf(test.expectedError, fileUnderTest))
			}
		})
	}
}

func fixturePath(fixtureName string) string {
	wd, _ := os.Getwd()
	return path.Join(wd, "../../fixtures/", fixtureName)
}
