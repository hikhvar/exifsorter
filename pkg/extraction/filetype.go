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
	"os"

	"github.com/h2non/filetype"
	"github.com/pkg/errors"
)

// IsVideoOrImage return true if the given file is a video or an image
func IsVideoOrImage(fname string) (bool, error) {
	// Open a file descriptor
	file, err := os.Open(fname)
	if err != nil {
		return false, errors.Wrap(err, "could not open file to determine file type")
	}
	defer file.Close()

	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	_, err = file.Read(head)
	if err != nil {
		return false, errors.Wrap(err, "could not read file header to determine file type")
	}
	return filetype.IsImage(head) || filetype.IsVideo(head), nil
}
