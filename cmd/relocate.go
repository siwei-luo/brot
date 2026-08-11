package cmd

import (
	"github.com/siwei-luo/brot/pkg"
	"github.com/spf13/cobra"
)

var dryRunRelocate bool = false

var relocateCmd = &cobra.Command{
	Use:   "relocate",
	Short: "Rule based file move/copy",
	Long: `Define custom rules to move/copy files around.

relocate:
  - name: example rule
    src: $HOME/Downloads
    dst: $HOME/Documents
    patterns:
      - "*.pdf"
    mode: move
	`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Relocate(dryRunRelocate)
	},
}

func init() {
	rootCmd.AddCommand(relocateCmd)

	relocateCmd.Flags().BoolVarP(&dryRunRelocate, "dry-run", "d", false, "Do not actually move/copy anything.")
}
