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
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestCaptureDate(t *testing.T) {
	tests := []struct {
		name          string
		timeStamp     time.Time
		expectedError string
	}{
		{
			name:      "sample1.JPG",
			timeStamp: parseTimeString(t, "2015-12-24 13:59:17 +0100 CET"),
		},
		{
			name: "sample2.mp4",
		},
		{
			name: "sample3.txt",
		},
		{
			name:          "sample-not-exist",
			expectedError: "could not open file to examine capture date: open /home/christoph/workspace/GO/src/github.com/hikhvar/exifsorter/fixtures/sample-not-exist: no such file or directory",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts, err := CaptureDate(fixturePath(test.name))
			assert.Equal(t, test.timeStamp, ts)
			if test.expectedError != "" {
				assert.EqualError(t, err, test.expectedError)
			}
		})
	}
}

func parseTimeString(t *testing.T, ts string) time.Time {
	ti, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", ts)
	if err != nil {
		t.Fatalf("broken test setup: %s", err.Error())
	}
	return ti
}
