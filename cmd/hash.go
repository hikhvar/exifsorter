package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/hikhvar/exifsorter/pkg/extraction"

	"github.com/hikhvar/exifsorter/pkg/exploration"
	"github.com/spf13/cobra"
	"github.com/timshannon/bolthold"
)

type HashedFile struct {
	Filepath    string
	Hash        string
	CaptureDate time.Time
}

// listCmd represents the list command
var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Hash all images in the directory, result will be written into a boltdb",
	Long:  `Hash all images in the directory, result will be written into a boltdb`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		inputDir := cmd.Flag("directory").Value.String()
		_, files, err := exploration.InitialFiles(inputDir, nil)
		if err != nil {
			fmt.Printf("could not list all files %s", err.Error())
			os.Exit(1)
		}

		store, err := bolthold.Open(cmd.Flag("database").Value.String(), 0666, nil)
		if err != nil {
			fmt.Printf("failed to open hash database: %v\n", err)
			os.Exit(1)
		}
		defer store.Close()

		for i, fp := range files {
			fmt.Printf("Hash image %d of %d\n", i+1, len(files))
			hf, err := extractFileInfo(inputDir, fp)
			if err != nil {
				fmt.Printf("failed to get image data: %v \n", err)
				continue
			}
			if hf.Filepath == "" {
				continue
			}
			err = store.Insert(hf.Filepath, hf)
			if err != nil {
				fmt.Printf("failed to store data in bolddb: %v \n", err)
			}
		}
	},
}

func extractFileInfo(directory string, fp string) (HashedFile, error) {
	mf, err := extraction.ReadFile(fp)
	if err != nil {
		return HashedFile{}, errors.Wrap(err, "failed to read source file")
	}
	isImage, err := mf.IsImage()
	if err != nil {
		return HashedFile{}, errors.Wrap(err, "not a video or image")
	}
	if !isImage {
		return HashedFile{}, nil
	}
	hash, err := extraction.HashImage(fp)
	if err != nil {
		return HashedFile{}, errors.Wrap(err, "failed to hash image")
	}
	cd, err := extraction.CaptureDate(fp)
	if err != nil {
		return HashedFile{}, errors.Wrap(err, "failed to extract capture date")

	}
	relPath, err := filepath.Rel(directory, fp)
	if err != nil {
		return HashedFile{}, errors.Wrap(err, "failed to compute relative path")
	}
	hf := HashedFile{
		Filepath:    relPath,
		Hash:        hash.ToString(),
		CaptureDate: cd,
	}
	return hf, nil
}

func init() {
	rootCmd.AddCommand(hashCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	hashCmd.PersistentFlags().StringP("directory", "d", "", "directory to hash")
	hashCmd.PersistentFlags().StringP("database", "o", "hash.db", "boltdb file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
