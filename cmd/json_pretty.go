package cmd

import (
    "fmt"

    "dt/internal/cliio"
    "dt/internal/jsonutil"
    "github.com/spf13/cobra"
)

func init() {
    jsonCmd.AddCommand(jsonPrettyCmd)
    rootCmd.AddCommand(jsonCmd)
}

var jsonCmd = &cobra.Command{
    Use:   "json",
    Short: "JSON utilities",
}

var prettyIndent int

var jsonPrettyCmd = &cobra.Command{
    Use:   "pretty",
    Short: "Pretty-print JSON (handles stringified JSON)",
    Example: `echo '"{\"a\":1,\"b\":[1,2]}"' | dt json pretty
dt json pretty '{"a":1}' --indent 2`,
    RunE: func(cmd *cobra.Command, args []string) error {
        in, err := cliio.ReadAll(args)
        if err != nil {
            return err
        }
        if prettyIndent < 0 {
            prettyIndent = 2
        }
        out, err := jsonutil.Pretty(in, prettyIndent)
        if err != nil {
            return err
        }
        fmt.Println(string(out))
        return nil
    },
}

func init() { // flags init
    jsonPrettyCmd.Flags().IntVar(&prettyIndent, "indent", 2, "number of spaces to indent")
    // shell completion for indent small set
    jsonPrettyCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
        opts := []string{"2\tdefault", "4", "0\tcompact"}
        return opts, cobra.ShellCompDirectiveNoFileComp
    }
}
