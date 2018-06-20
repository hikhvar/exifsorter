package files

import (
	"os"

	"syscall"

	"io"

	"path"

	"github.com/pkg/errors"
)

// IsNormalFile returns true if the given file is not a directory
func IsNormalFile(fname string) (bool, error) {
	fInfo, err := os.Stat(fname)
	if err != nil {
		return false, err
	}
	return !fInfo.IsDir(), nil
}

// File copies src file to dst. dst is truncated or created if not present. The FileMode and Modtimes are preserved.
func Copy(src, dst string) error {
	fInfo, err := os.Stat(src)
	if err != nil {
		return errors.Wrap(err, "can not get file info of src")
	}
	if fInfo.IsDir() {
		return errors.New("src is a directory")
	}
	targetDiskSize, err := getFreeDiskSize(dst)
	if err != nil {
		return errors.Wrap(err, "can not get remaining disk size in dst")
	}
	if targetDiskSize < uint64(fInfo.Size()) {
		return errors.New("not enough space left in dst")
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, "can not open src file")
	}
	defer srcFile.Close()
	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_TRUNC|os.O_CREATE, fInfo.Mode())
	if err != nil {
		return errors.Wrap(err, "can not open dst file")
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return errors.Wrap(err, "error while copying file")
	}
	err = os.Chtimes(dst, fInfo.ModTime(), fInfo.ModTime())
	if err != nil {
		return errors.Wrap(err, "can not copy change times from src")
	}
	return dstFile.Sync()
}

// getFreeDiskSize returns the available disk size in bytes
func getFreeDiskSize(dir string) (uint64, error) {
	var stat syscall.Statfs_t
	fInfo, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			dir = path.Dir(dir)
		} else {
			return 0, errors.Wrap(err, "can not get file info of dir")
		}
	} else if !fInfo.IsDir() {
		dir = path.Dir(dir)
	}
	err = syscall.Statfs(dir, &stat)
	if err != nil {
		return 0, errors.Wrap(err, "failed syscall Statfs")
	}

	// Available blocks * size per block = available space in bytes
	return stat.Bavail * uint64(stat.Bsize), nil
}
