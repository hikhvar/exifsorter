package archive

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/hikhvar/exifsorter/pkg/extraction"
)

func NewOSFileSystem() FileSystem {
	return FileSystem{
		fd:            os.Remove,
		linker:        os.Link,
		mkdir:         os.MkdirAll,
		isMedia:       extraction.IsVideoOrImage,
		dateExtractor: extraction.CaptureDate,
	}
}

func NewLoggingFileSystem() FileSystem {
	return FileSystem{
		fd: func(file string) error {
			log.Printf("[DRY-RUN] will delete file: %s", file)
			return nil
		},
		linker: func(old, new string) error {
			log.Printf("[DRY-RUN] link %s to %s", old, new)
			return nil
		},
		mkdir: func(dirPath string, perm os.FileMode) error {
			log.Printf("[DRY-RUN] create directory %s with mode %s", dirPath, perm)
			return nil
		},
	}
}

type FileSystem struct {
	fd            FileDeleter
	linker        Linker
	mkdir         DirectoryCreator
	isMedia       IsMedia
	dateExtractor DateExtractor
}

// EnsureAbsent removes the given directory and returns an error if file is not deleted
func (fs FileSystem) EnsureAbsent(name string) error {
	err := fs.fd(name)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing file before creating link: %w", err)
	}
	return nil
}

// EnsureDirectory creates the directory recursive
func (fs FileSystem) EnsureDirectory(name string) error {
	return fs.mkdir(name, os.ModePerm)
}

// createLinks create a symlink from every path in paths to the given target
func (fs FileSystem) CreateLinks(paths []string, target string) error {
	for _, p := range paths {
		err := fs.EnsureAbsent(p)
		if err != nil {
			return errors.Wrap(err, "can't ensure file is not currently absent")
		}
		err = fs.EnsureDirectory(filepath.Dir(p))
		if err != nil {
			return errors.Wrap(err, "can not create directory for link")
		}
		err = fs.linker(target, p)
		if err != nil {
			return errors.Wrap(err, "can not hard link to all archive")
		}
	}
	return nil
}
