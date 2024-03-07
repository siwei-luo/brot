/*
Copyright © 2021-2023 Siwei Luo <siwei@lu0.org>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/siwei-luo/brot/pkg"
	"github.com/spf13/cobra"
)

var dryRunCleanup bool = false

// cleanupCmd represents the cleanup command
var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Rule based file cleanup",
	Long: `Define custom rules to cleanup obsolete files.

cleanup:
  - name: mac os foo
    src: $HOME/Downloads
    patterns:
      - ".DS_Store"
      - ".AppleDouble"
      - ".LSOverride"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Cleanup(dryRunCleanup)
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cleanupCmd.Flags().BoolVarP(&dryRunCleanup, "dry-run", "d", false, "Do not actually delete anything.")
}
