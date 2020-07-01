package archive

import (
	"path/filepath"
	"time"

	"os"

	"fmt"
	"path"

	"crypto/sha256"
	"hash"

	"github.com/fsnotify/fsnotify"
	"github.com/hikhvar/exifsorter/pkg/extraction"
	"github.com/hikhvar/exifsorter/pkg/files"
	"github.com/pkg/errors"
)

const targetTimeFormat = "20060102_150405"

type Watcher interface {
	Channels() (chan fsnotify.Event, chan error)
}

type Copier func(src, dst string, hFunc hash.Hash) (hashSum []byte, err error)
type Linker func(old, new string) error
type DirectoryCreator func(dirPath string, perm os.FileMode) error
type DateExtractor func(fname string) (time.Time, error)
type IsMedia func(fname string) (bool, error)

type Algorithm struct {
	archiveDir       string
	sourceDir        string
	copier           Copier
	directoryCreator DirectoryCreator
	extractor        DateExtractor
	isMedia          IsMedia
}

// NewArchive returns a new Algorithm.
func NewAlgorithm(src, dst string) *Algorithm {
	return &Algorithm{
		archiveDir:       dst,
		sourceDir:        src,
		copier:           files.Copy,
		directoryCreator: os.MkdirAll,
		extractor:        extraction.CaptureDate,
		isMedia:          extraction.IsVideoOrImage,
	}
}

// Init creates all required target directories
func (a *Algorithm) Init() error {
	err := a.directoryCreator(a.allArchiveDir(), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "could not create target dir '%s'", a.allArchiveDir())
	}

	err = a.directoryCreator(a.originArchiveDir(), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "could not create target dir '%s'", a.originArchiveDir())
	}
	return nil
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

	tmpFile := path.Join(targetDir, "exifsorter.tmp")
	sum, err := a.copier(fname, tmpFile, sha256.New224())
	if err != nil {
		return tmpFile, errors.Wrap(err, "could not copy file and compute checksum")
	}

	targetFileName := fmt.Sprintf("%s_%s%s", date.Format(targetTimeFormat), fmt.Sprintf("%x", sum)[0:8], path.Ext(fname))
	targetFilePath := path.Join(targetDir, targetFileName)
	err = os.Rename(tmpFile, targetFilePath)
	if err != nil {
		return tmpFile, errors.Wrap(err, "could not mv temporary file to target name")
	}

	allArchiveName := path.Join(a.allArchiveDir(), targetFileName)
	originArchiveName, err := a.originArchiveFileName(fname, targetFileName)
	if err != nil {
		return targetFilePath, errors.Wrap(err, "failed to determine relative path")
	}
	return targetFilePath, a.createLinks([]string{allArchiveName, originArchiveName}, targetFilePath)
}

// createLinks create a symlink from every path in paths to the given target
func (a *Algorithm) createLinks(paths []string, target string) error {
	for _, p := range paths {
		os.Remove(p) // nolint:errcheck
		err := a.directoryCreator(filepath.Dir(p), os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "can not create directory for link")
		}
		err = os.Link(target, p)
		if err != nil {
			return errors.Wrap(err, "can not hard link to all archive")
		}
	}
	return nil
}

func (a *Algorithm) originArchiveFileName(sourceFileName string, targetFileName string) (string, error) {
	pathInSrc, err := filepath.Rel(a.sourceDir, sourceFileName)
	if err != nil {
		return "", errors.Wrap(err, "failed to compute relative path in source")
	}
	dirName := filepath.Dir(pathInSrc)
	pathInOrigin := path.Join(dirName, targetFileName)
	return path.Join(a.originArchiveDir(), pathInOrigin), nil
}

func (a *Algorithm) allArchiveDir() string {
	return path.Join(a.archiveDir, "all")
}

func (a *Algorithm) originArchiveDir() string {
	return path.Join(a.archiveDir, "origin")
}

func getYearMonth(t time.Time) (int, int) {
	return t.Year(), int(t.Month())
}
