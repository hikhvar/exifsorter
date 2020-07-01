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
	"context"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/hikhvar/exifsorter/pkg/archive"
	"github.com/hikhvar/exifsorter/pkg/exploration"
	"github.com/hikhvar/exifsorter/pkg/files"
	"github.com/spf13/cobra"
)

var ignorePatterns []string

// sortCmd represents the sort command
var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "sorts media data according to their exif metadata",
	Long:  `sorts media data according to their exif metadata`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancelFunc := context.WithCancel(context.Background())
		defer cancelFunc()
		srcDir, dstDir := srcAndDstDir(cmd)
		a := archive.NewAlgorithm(srcDir, dstDir)
		err := a.Init()
		if err != nil {
			fmt.Printf("failed to create target directories: %v", err)
			os.Exit(1)
		}

		ignores, err := exploration.GobwasMatcherFromPatterns(ignorePatterns)
		if err != nil {
			fmt.Printf("not valid globs '%v': %v", ignorePatterns, err.Error())
			os.Exit(1)
		}
		dirs, fs, err := exploration.InitialFiles(srcDir, ignores)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, f := range fs {
			n, err := a.Sort(f)
			if err != nil && err.Error() != "given file is not a media file" {
				fmt.Printf("Can't sort file %v: %v", f, err.Error())
			} else {
				fmt.Printf("%s\t-->\t%s\n", f, n)
			}
		}
		fmt.Println("finished intial run. Watch folder for changes.")
		watcher, err := exploration.NewRecursiveWatcher(ctx, ignores, dirs...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for {
			select {
			case err = <-watcher.Errors:
				fmt.Println(err)
			case e := <-watcher.Events:
				if e.Op == fsnotify.Remove {
					break
				}
				f := e.Name
				normalFile, err := files.IsNormalFile(f)
				if err == nil {
					if normalFile {
						n, err := a.Sort(f)
						if err != nil && err.Error() != "given file is not a media file" {
							fmt.Printf("%v: %v", f, err.Error())
						} else {
							fmt.Printf("%s\t-->\t%s\n", f, n)
						}

					}
				} else {
					fmt.Printf("could not stat file: %v\n", err)
				}
			}
		}
	},
}

func srcAndDstDir(cmd *cobra.Command) (string, string) {
	return cmd.Flag("source").Value.String(), cmd.Flag("target").Value.String()
}

func init() {
	rootCmd.AddCommand(sortCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	sortCmd.PersistentFlags().StringP("source", "s", "", "source directory")

	sortCmd.PersistentFlags().StringP("target", "t", "", "target directory")

	sortCmd.PersistentFlags().StringArrayVarP(&ignorePatterns, "ignores", "i", []string{"**.@__thumb**", "**.syncthing.*tmp", "**.!sync"}, "file patterns to ignore. For supported patterns see https://github.com/gobwas/glob .")

	sortCmd.PersistentFlags().BoolP("dry-run", "d", false, "dry run. Don't edit anything.")
}
