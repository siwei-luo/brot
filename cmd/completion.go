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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// completionCmd represents the completion command
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
