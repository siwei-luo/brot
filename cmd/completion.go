package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|fish|powershell|zsh]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(brot completion bash)

# To load completions for each session, execute once:
Linux:
  $ brot completion bash > /etc/bash_completion.d/brot
MacOS:
  $ brot completion bash > /usr/local/etc/bash_completion.d/brot

Fish:

$ brot completion fish | source

# To load completions for each session, execute once:
$ brot completion fish > ~/.config/fish/completions/brot.fish

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ brot completion zsh > "${fpath[1]}/_brot"

# You will need to start a new shell for this setup to take effect.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			if err := cmd.Root().GenBashCompletion(os.Stdout); err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Fatal("error")
			}
		case "fish":
			if err := cmd.Root().GenFishCompletion(os.Stdout, true); err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Fatal("error")
			}
		case "powershell":
			if err := cmd.Root().GenPowerShellCompletion(os.Stdout); err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Fatal("error")
			}
		case "zsh":
			if err := cmd.Root().GenZshCompletion(os.Stdout); err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Fatal("error")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
