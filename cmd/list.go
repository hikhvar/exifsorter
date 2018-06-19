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

	"github.com/hikhvar/exifsorter/pkg/exploration"
	"github.com/hikhvar/exifsorter/pkg/extraction"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the found exif meta data for a subdirectory",
	Long:  `List the found exif meta data for a subdirectory`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, files, err := exploration.InitialFiles(args[0])
		if err != nil {
			fmt.Printf("could not list all files %s", err.Error())
		}
		for _, f := range files {
			voi, err := extraction.IsVideoOrImage(f)
			if err != nil {
				fmt.Printf("not a video or image %s: %s\n", f, err.Error())
			} else if voi {
				date, err := extraction.CaptureDate(f)
				if err != nil {
					fmt.Printf("could not determine capture date %s: %s\n", f, err.Error())
				} else {
					fmt.Printf("exif date of file %s is: %v\n", f, date)
				}

			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	listCmd.PersistentFlags().StringP("directory", "d", "", "directory to list")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
