package cmd

import (
	"github.com/siwei-luo/brot/pkg"
	"github.com/spf13/cobra"
)

var dryRunCleanup bool = false

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

	cleanupCmd.Flags().BoolVarP(&dryRunCleanup, "dry-run", "d", false, "Do not actually delete anything.")
}
