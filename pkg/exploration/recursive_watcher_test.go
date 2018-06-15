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
	"context"
	"os"
	"testing"

	"io/ioutil"
	"time"

	"path"

	"github.com/pkg/errors"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
)

type touchFile struct {
	name  string
	isDir bool
}

func TestNewRecursiveWatcher(t *testing.T) {
	tests := []struct {
		name           string
		dir            string
		cleanDir       bool
		filesToTouch   []touchFile
		expectedEvents []fsnotify.Event
		expectedError  error
		expectFailure  bool
	}{
		{
			name:     "emptyDir",
			dir:      createTempDir(t),
			cleanDir: true,
		},
		{
			name:     "emtpy Dir with supdirectory",
			dir:      createTempDir(t),
			cleanDir: true,
			filesToTouch: []touchFile{
				{
					isDir: true,
					name:  "foo",
				},
			},
			expectedEvents: []fsnotify.Event{
				{
					Op:   fsnotify.Create,
					Name: "foo",
				},
			},
		},
		{
			name:     "emtpy Dir with create in supdirectory",
			dir:      createTempDir(t),
			cleanDir: true,
			filesToTouch: []touchFile{
				{
					isDir: true,
					name:  "foo",
				},
				{
					isDir: false,
					name:  "foo/bar",
				},
			},
			expectedEvents: []fsnotify.Event{
				{
					Op:   fsnotify.Create,
					Name: "foo",
				},
				{
					Op:   fsnotify.Create,
					Name: "foo/bar",
				},
			},
		},
		{
			name:          "dir does not exists",
			dir:           "/tmp/foo-bar",
			cleanDir:      false,
			expectFailure: true,
			expectedError: errors.New("could not add /tmp/foo-bar to watcher: no such file or directory"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.cleanDir {
				defer os.RemoveAll(test.dir)
			}
			ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancelFunc()
			w, err := NewRecursiveWatcher(ctx, test.dir)
			if test.expectedError != nil {
				if !assert.NotNil(t, err) {
					return
				}
				assert.Equal(t, test.expectedError.Error(), err.Error())
			}
			if test.expectFailure {
				assert.Nil(t, w)
				return
			}
			receivedEvents := make([]fsnotify.Event, 0)
			touchFiles(t, test.dir, test.filesToTouch)
			for {
				select {
				case <-ctx.Done():
					goto END
				case e := <-w.Events:
					receivedEvents = append(receivedEvents, e)
				}
			}
		END:
			joinExpectedEventsWithDir(test.dir, test.expectedEvents)
			assert.ElementsMatch(t, test.expectedEvents, receivedEvents)
		})
	}
}

func joinExpectedEventsWithDir(testDir string, expectedEvents []fsnotify.Event) {
	for i := range expectedEvents {
		expectedEvents[i].Name = path.Join(testDir, expectedEvents[i].Name)
	}
}

func createTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "TestNewRecursiveWatcher")
	if err != nil {
		t.Fatalf("broken test setup. can not create tempDir %s", err.Error())
	}
	return dir
}

func touchFiles(t *testing.T, root string, files []touchFile) {

	for _, fn := range files {
		time.Sleep(100 * time.Millisecond)
		name := path.Join(root, fn.name)
		if fn.isDir {
			err := os.MkdirAll(name, os.ModePerm)
			if err != nil {
				t.Fatalf("broken test setup: %s", err.Error())
			}
		} else {
			f, err := os.Create(name)
			if err != nil {
				t.Fatalf("broken test setup: %s", err.Error())
			}
			f.Close()
		}

	}
}
