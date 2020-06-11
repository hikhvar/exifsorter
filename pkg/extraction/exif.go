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
	"time"

	"os"

	"github.com/pkg/errors"

	"github.com/xor-gate/goexif2/exif"
	"github.com/xor-gate/goexif2/mknote"
)

const (
	//	dateAndTime          = "Date and Time"
	//	dateAndTimeDigitized = "Date and Time (Digitized)"
	//	dateAndTimeOriginial = "Date and Time (Original)"
	//	timeFormat       = "2006:01:02 15:04:05"
	noInfoFoundError = "could neither read exif meta data nor file modification time"
)

func init() {
	exif.RegisterParsers(mknote.All...)
}

// CaptureDate returns the point in time the capturing device created the media file
func CaptureDate(fname string) (time.Time, error) {
	fInfo, fInfoErr := os.Stat(fname)
	f, err := os.Open(fname)
	if err != nil {
		if fInfoErr == nil {
			return fInfo.ModTime(), nil
		}
		return time.Time{}, errors.Wrap(err, "failed to open or fstat file.")
	}
	x, err := exif.Decode(f)
	//data, err := exif.Read(fname)
	if err != nil {
		if fInfoErr == nil {
			return fInfo.ModTime(), nil
		}
		return time.Time{}, errors.Wrap(err, noInfoFoundError)
	}
	tm, err := x.DateTime()
	if err != nil {
		if fInfoErr == nil {
			return fInfo.ModTime(), nil
		}
		return time.Time{}, errors.Wrap(err, noInfoFoundError)
	}
	return tm, nil
}
