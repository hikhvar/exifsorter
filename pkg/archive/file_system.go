package archive

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/hikhvar/exifsorter/pkg/extraction"
)

func NewOSFileSystem() FileSystem {
	return FileSystem{
		fd:            os.Remove,
		linker:        os.Link,
		mkdir:         os.MkdirAll,
		stater:        os.Stat,
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
		stater: func(name string) (os.FileInfo, error) {
			log.Printf("[DRY-RUN] stat %s", name)
			return FakeFileInfo{name}, nil
		},
	}
}

type FileSystem struct {
	fd            FileDeleter
	linker        Linker
	stater        Stater
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

func (fs FileSystem) EqualSize(oldFile, newFile string) (bool, error) {
	oldStats, err := fs.stater(oldFile)
	if err != nil {
		return false, fmt.Errorf("failed to stat old file: %w", err)
	}
	newStats, err := fs.stater(newFile)
	if err != nil {
		return false, fmt.Errorf("failed to stat new file: %w", err)
	}
	return oldStats.Size() == newStats.Size(), nil
}

type FakeFileInfo struct {
	name string
}

func (f FakeFileInfo) Name() string {
	return f.name
}

func (f FakeFileInfo) Size() int64 {
	return 0
}

func (f FakeFileInfo) Mode() fs.FileMode {
	//TODO implement me
	panic("implement me")
}

func (f FakeFileInfo) ModTime() time.Time {
	//TODO implement me
	panic("implement me")
}

func (f FakeFileInfo) IsDir() bool {
	//TODO implement me
	panic("implement me")
}

func (f FakeFileInfo) Sys() any {
	//TODO implement me
	panic("implement me")
}
