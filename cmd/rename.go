package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/siwei-luo/brot/pkg"
	"github.com/spf13/cobra"
)

var renameChecksum string
var renamePattern string
var renamePrefixLength int

var renameCmd = &cobra.Command{
	Use:   "rename [directory]",
	Short: "Rename files based on checksum",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]

		if renameChecksum != "sha256" {
			log.Fatal("only sha256 is supported for --checksum")
		}

		if err := pkg.RenameFiles(dir, renamePattern, renamePrefixLength); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
	renameCmd.Flags().StringVar(&renameChecksum, "checksum", "sha256", "checksum algorithm to use")
	renameCmd.Flags().StringVar(&renamePattern, "pattern", "*.jpg", "file pattern to match")
	renameCmd.Flags().IntVar(&renamePrefixLength, "prefix-length", 5, "length of the numeric prefix (zero-padded)")
}
