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

package cmd

import (
	"fmt"
	"os"

	"github.com/hikhvar/exifsorter/pkg/archive"
	"github.com/hikhvar/exifsorter/pkg/exploration"
	"github.com/spf13/cobra"
)

// sortCmd represents the sort command
var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "sorts media data according to their exif metadata",
	Long:  `sorts media data according to their exif metadata`,
	Run: func(cmd *cobra.Command, args []string) {
		a := archive.NewAlgorithm(cmd.Flag("target").Value.String())
		_, files, err := exploration.InitialFiles(cmd.Flag("source").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, f := range files {
			fmt.Println(f)
			err = a.Sort(f)
			if err != nil && err.Error() != "given file is not a media file" {
				fmt.Printf("%v: %v", f, err.Error())
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(sortCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	sortCmd.PersistentFlags().StringP("source", "s", "", "source directory")

	sortCmd.PersistentFlags().StringP("target", "t", "", "target directory")

	sortCmd.PersistentFlags().BoolP("dry-run", "d", false, "dry run. Don't edit anything.")
}
