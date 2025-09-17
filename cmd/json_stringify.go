package cmd

import (
    "fmt"

    "dt/internal/cliio"
    "dt/internal/jsonutil"
    "github.com/spf13/cobra"
)

var (
    stringifyCompact  bool
    stringifyNoQuotes bool
)

func init() {
    jsonCmd.AddCommand(jsonStringifyCmd)
}

var jsonStringifyCmd = &cobra.Command{
    Use:   "stringify",
    Short: "Convert JSON to a single JSON string (escaped)",
    Example: `cat obj.json | dt json stringify --compact
dt json stringify '{"a":1}' --no-quotes`,
    RunE: func(cmd *cobra.Command, args []string) error {
        in, err := cliio.ReadAll(args)
        if err != nil {
            return err
        }
        out, err := jsonutil.Stringify(in, stringifyCompact, stringifyNoQuotes)
        if err != nil {
            return err
        }
        fmt.Println(string(out))
        return nil
    },
}

func init() {
    jsonStringifyCmd.Flags().BoolVar(&stringifyCompact, "compact", false, "minify before stringifying")
    jsonStringifyCmd.Flags().BoolVar(&stringifyNoQuotes, "no-quotes", false, "omit surrounding quotes")
}

