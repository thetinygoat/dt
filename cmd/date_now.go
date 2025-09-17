package cmd

import (
    "fmt"
    "time"

    "dt/internal/dateutil"
    "github.com/spf13/cobra"
)

var (
    dateFormat string
    dateLayout string
    dateUTC    bool
)

func init() {
    rootCmd.AddCommand(dateCmd)
    dateCmd.AddCommand(dateNowCmd)
    dateCmd.AddCommand(dateToEpochCmd)
    dateCmd.AddCommand(dateFromEpochCmd)
    dateCmd.AddCommand(dateAddCmd)
}

var dateCmd = &cobra.Command{Use: "date", Short: "Date and time helpers"}

var dateNowCmd = &cobra.Command{
    Use:   "now",
    Short: "Print current time",
    RunE: func(cmd *cobra.Command, args []string) error {
        t := time.Now()
        fmt.Println(dateutil.FormatTime(t, dateFormat, dateLayout, dateUTC))
        return nil
    },
}

func init() {
    // shared flags on date namespace where useful
    for _, c := range []*cobra.Command{dateNowCmd} {
        c.Flags().StringVar(&dateFormat, "format", "rfc3339", "output format: rfc3339|unix|unixms|layout|<Go layout>")
        c.Flags().StringVar(&dateLayout, "layout", "", "when --format=layout, Go time layout to use")
        c.Flags().BoolVar(&dateUTC, "utc", false, "print in UTC")
    }
    // completion for --format
    dateNowCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
        return []string{"rfc3339", "unix", "unixms", "layout"}, cobra.ShellCompDirectiveNoFileComp
    })
}

