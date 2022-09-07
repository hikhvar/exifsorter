package extraction

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/pkg/errors"

	"github.com/corona10/goimagehash"
)

func HashImage(fname string) (*goimagehash.ExtImageHash, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open image file")
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode image")
	}
	hash, err := goimagehash.ExtPerceptionHash(img, 8, 8)
	return hash, errors.Wrap(err, "failed to hash image")
}
