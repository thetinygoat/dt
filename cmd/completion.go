package cmd

import (
    "os"

    "github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(completionCmd)
}

var completionCmd = &cobra.Command{
    Use:   "completion [bash|zsh|fish|powershell]",
    Short: "Generate shell completion scripts",
    Long: `To enable completions:

Bash:
  source <(dt completion bash)
  # or
  dt completion bash > /etc/bash_completion.d/dt

Zsh:
  dt completion zsh > "${fpath[1]}/_dt"
  # ensure 'compinit' is run in your zshrc

Fish:
  dt completion fish | source
  # or
  dt completion fish > ~/.config/fish/completions/dt.fish

PowerShell:
  dt completion powershell | Out-String | Invoke-Expression
`,
    Args: cobra.ExactValidArgs(1),
    ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
    RunE: func(cmd *cobra.Command, args []string) error {
        switch args[0] {
        case "bash":
            return rootCmd.GenBashCompletion(os.Stdout)
        case "zsh":
            return rootCmd.GenZshCompletion(os.Stdout)
        case "fish":
            return rootCmd.GenFishCompletion(os.Stdout, true)
        case "powershell":
            return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
        default:
            return nil
        }
    },
}

