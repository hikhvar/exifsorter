package archive

import (
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type FileDeleter func(file string) error

type DeDupTask struct {
	// ToKeep is the file path of the original to keep
	ToKeep string
	// ReCreateLinks are files which should be hard link to ToKeep. If they are already present, they should be deleted and recreated
	ReCreateLinks []string
	// DeleteFiles are files
	DeleteFiles []string
}

// DeduplicateAll deduplicates all given files in the directory. This method actually executes the file operations if noDryRun is set.
func DeduplicateAll(archiveRoot string, duplicates [][]string, creator FileSystem) error {

	for _, duplicateFiles := range duplicates {
		task, err := DeDuplicate(archiveRoot, duplicateFiles)
		if err != nil {
			return fmt.Errorf("failed to compute deduplicateTask for %s: %w", duplicateFiles, err)
		}
		err = creator.CreateLinks(task.ReCreateLinks, task.ToKeep)
		if err != nil {
			return fmt.Errorf("failed to create links to: %w", err)
		}
		for _, toDelete := range task.DeleteFiles {
			err := creator.EnsureAbsent(toDelete)
			if err != nil {
				return fmt.Errorf("failed to delete file: %w", err)
			}
		}
	}
	return nil
}

// DeDuplicate files in the given archiveRoot. All files in duplicateFiles must start with the prefix archiveRoot.
// This function assumes the canonical archive layout:
// /archiveRoot/
//
//	/ YEAR
//	   / Month1
//	   / Month2
//	/ origin
//	   / dirOne
//	   / dirTwo
//
// The file in DedupTask.ToKeep will be in the directory /YEAR/MONTH. If there are multiple files in the /YEAR/MONTH directories, the first file is kept.
// At most one file in every directory below /origin is kept.
func DeDuplicate(archiveRoot string, duplicateFiles []string) (DeDupTask, error) {
	sort.Strings(duplicateFiles)
	ret := DeDupTask{}
	foundInDirectory := make(map[string]struct{})
	for _, f := range duplicateFiles {
		inArchive, err := pathInArchive(archiveRoot, f)
		if err != nil {
			return DeDupTask{}, fmt.Errorf("failed to find path in directory: %w", err)
		}
		if isCalendarStoredFile(inArchive) {
			if ret.ToKeep == "" {
				ret.ToKeep = f
			} else {
				ret.DeleteFiles = append(ret.DeleteFiles, f)
			}
			continue
		}
		dir := filepath.Dir(inArchive)
		if _, found := foundInDirectory[dir]; found {
			ret.DeleteFiles = append(ret.DeleteFiles, f)
			continue
		}
		foundInDirectory[dir] = struct{}{}
		ret.ReCreateLinks = append(ret.ReCreateLinks, f)
	}
	if ret.ToKeep == "" {
		return DeDupTask{}, fmt.Errorf("there is no file in calendar directory")
	}
	return ret, nil
}

// pathInArchive returns the relative path within the archive. Returns an error if the file is not within the archiveRoot
func pathInArchive(archiveRoot string, filename string) (string, error) {
	rel, err := filepath.Rel(archiveRoot, filename)
	if strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("path is outside of archive")
	}
	return rel, err
}

// isCalendarStoredFile returns true if the file is stored in a calendar directory within the archive. The filename must be a relative path within the archive.
func isCalendarStoredFile(filename string) bool {
	matched, err := path.Match("[0-9][0-9][0-9][0-9]/[0-9][0-9]/*", filename)
	if err != nil {
		panic(err)
	}
	return matched
}
