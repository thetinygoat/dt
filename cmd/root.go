package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "dt",
    Short: "dt: day-to-day developer toolbox",
    Long:  "dt is a small, focused CLI to speed up daily dev tasks (JSON, dates, base64, UUIDs, env conversions).",
}

// Execute is the program entry from main.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

