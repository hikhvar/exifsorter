package cmd

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hikhvar/exifsorter/pkg/archive"
)

const (
	directoryParameterName = "directory"
	inputParameterName     = "input"
	delimiterParameterName = "delimiter"
	dryrunParameterName    = "dry-run"
)

// dedupCmd represents the dedup command
var dedupCmd = &cobra.Command{
	Use:   "dedup",
	Short: "Deduplicate the files in the given directory",
	Long:  `Deduplicate the files in the given directory. The duplicated files must be given in a file. The format of the given file is the output of: https://gitlab.com/opennota/findimagedupes`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		archiveRoot := cmd.Flag(directoryParameterName).Value.String()
		inputFilePath := cmd.Flag(inputParameterName).Value.String()
		delimiter := cmd.Flag("delimiter").Value.String()

		f, err := os.Open(inputFilePath)
		if err != nil {
			log.Printf("can't open input file: %s", err)
			os.Exit(1)
		}
		if len(delimiter) > 1 {
			log.Printf("can only use a single character as delimiter. '%s' has the length %d", delimiter, len(delimiter))
			os.Exit(1)
		}
		if len(delimiter) < 1 {
			log.Printf("Empty string not allowed as delimiter")
			os.Exit(1)
		}

		duplicates, err := readInput(f, delimiter)
		if err != nil {
			log.Printf("failed to read input file: %s", err)
			os.Exit(1)
		}

		dryRun, err := cmd.PersistentFlags().GetBool(dryrunParameterName)
		if err != nil {
			log.Printf("expected dry-run flag, didn't found it: %s", err)
		}

		var fs archive.FileSystem = archive.NewOSFileSystem()
		if dryRun {
			fs = archive.NewLoggingFileSystem()
		}
		err = archive.DeduplicateAll(archiveRoot, duplicates, fs)
		if err != nil {
			log.Printf("failed to deduplicate files: %s", err)
			os.Exit(1)
		}
	},
}

func readInput(reader io.Reader, delimiter string) ([][]string, error) {
	var ret [][]string
	s := bufio.NewScanner(reader)
	for s.Scan() {
		line := strings.Split(s.Text(), delimiter)
		ret = append(ret, line)
	}
	return ret, s.Err()
}

func init() {
	rootCmd.AddCommand(dedupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	dedupCmd.PersistentFlags().StringP(directoryParameterName, "", "", "directory to deduplicate in")
	dedupCmd.PersistentFlags().StringP(inputParameterName, "i", "", "path to a file with duplicated files")
	dedupCmd.PersistentFlags().StringP(delimiterParameterName, "", " ", "delimiter used in the file given by INPUT")
	dedupCmd.PersistentFlags().BoolP(dryrunParameterName, "", true, "don't deduplicate, only dry-run")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
