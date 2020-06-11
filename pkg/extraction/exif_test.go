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
			name:      "sample2.mp4",
			timeStamp: parseTimeString(t, "2016-04-02 09:23:56 +0200 CEST"),
		},
		{
			name:      "sample3.txt",
			timeStamp: parseTimeString(t, "2018-06-15 15:24:26.263360885 +0200 CEST"),
		},
		{
			name:          "sample-not-exist",
			expectedError: "failed to open or fstat file.: open %s: no such file or directory",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fileUnderTest := fixturePath(test.name)
			ts, err := CaptureDate(fileUnderTest)
			assert.Equal(t, test.timeStamp, ts)
			if test.expectedError != "" {
				assert.EqualError(t, err, fmt.Sprintf(test.expectedError, fileUnderTest))
			} else {
				assert.Nil(t, err)
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
