package archive

import (
	"time"

	"os"

	"fmt"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/hikhvar/exifsorter/pkg/extraction"
	"github.com/hikhvar/exifsorter/pkg/files"
	"github.com/pkg/errors"
)

type Watcher interface {
	Channels() (chan fsnotify.Event, chan error)
}

type Copier func(src, dst string) error
type DirectoryCreator func(dirPath string, perm os.FileMode) error
type DateExtractor func(fname string) (time.Time, error)
type IsMedia func(fname string) (bool, error)

type Algorithm struct {
	archiveDir       string
	copier           Copier
	directoryCreator DirectoryCreator
	extractor        DateExtractor
	isMedia          IsMedia
}

// NewArchive returns a new Algorithm.
func NewAlgorithm(dir string) *Algorithm {
	return &Algorithm{
		archiveDir:       dir,
		copier:           files.Copy,
		directoryCreator: os.MkdirAll,
		extractor:        extraction.CaptureDate,
		isMedia:          extraction.IsVideoOrImage,
	}
}

func (a *Algorithm) Sort(fname string) (string, error) {
	isMedia, err := a.isMedia(fname)
	if err != nil {
		return "", errors.Wrap(err, "could not determine media type")
	}
	if !isMedia {
		return "", errors.New("given file is not a media file")
	}
	date, err := a.extractor(fname)
	if err != nil {
		return "", errors.Wrap(err, "could not determine creation date of media file")
	}
	year, month := getYearMonth(date)
	targetDir := path.Join(a.archiveDir, fmt.Sprintf("%d/%02d", year, month))
	err = a.directoryCreator(targetDir, os.ModePerm)
	if err != nil {
		return "", errors.Wrapf(err, "could not create target dir '%s'", targetDir)
	}
	targetFile := path.Join(targetDir, path.Base(fname))
	return targetFile, a.copier(fname, targetFile)
}

func getYearMonth(t time.Time) (int, int) {
	return t.Year(), int(t.Month())
}
