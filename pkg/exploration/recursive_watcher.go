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

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

type RecursiveWatcher struct {
	watcher *fsnotify.Watcher
	ignores []Matcher
	Events  chan fsnotify.Event
	Errors  chan error
}

// NewRecursiveWatcher creates a new recursive file watcher. You can listen for errors and events via the channels
// Events and Errors
func NewRecursiveWatcher(ctx context.Context, ignores []Matcher, initialDirs ...string) (*RecursiveWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.Wrap(err, "could not create watcher")
	}
	for _, dir := range initialDirs {
		err = watcher.Add(dir)
		if err != nil {
			watcher.Close()
			return nil, errors.Wrapf(err, "could not add %s to watcher", dir)
		}
	}

	r := &RecursiveWatcher{
		watcher: watcher,
		ignores: ignores,
		Events:  make(chan fsnotify.Event, 10),
		Errors:  make(chan error),
	}
	go r.run(ctx)
	return r, nil
}

func (r *RecursiveWatcher) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			err := r.watcher.Close()
			if err != nil {
				r.Errors <- err
			}
			return
		case e := <-r.watcher.Errors:
			r.Errors <- e
		case e := <-r.watcher.Events:
			if !isIgnored(r.ignores, e.Name) {
				r.processEvent(e)
				r.Events <- e
			}
		}
	}
}

func (r *RecursiveWatcher) processEvent(e fsnotify.Event) {
	switch e.Op {
	case fsnotify.Create:
		finfo, err := os.Stat(e.Name)
		if err != nil {
			return
		}
		if finfo.IsDir() {
			r.watcher.Add(e.Name)
		}
	}
}
