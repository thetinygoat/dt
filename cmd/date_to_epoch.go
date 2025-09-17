package cmd

import (
    "fmt"
    "strings"

    "dt/internal/cliio"
    "dt/internal/dateutil"
    "github.com/spf13/cobra"
)

var (
    toEpochLayout string
    toEpochUTC    bool
    toEpochMs     bool
)

func init() {
    dateCmd.AddCommand(dateToEpochCmd)
}

var dateToEpochCmd = &cobra.Command{
    Use:   "to-epoch [times...]",
    Short: "Parse human time to Unix epoch",
    Long:  "Parses common human-readable time formats to Unix epoch seconds by default.",
    RunE: func(cmd *cobra.Command, args []string) error {
        var in string
        if cliio.IsInputFromPipe() {
            b, err := cliio.ReadAll(nil)
            if err != nil {
                return err
            }
            in = string(b)
        } else {
            in = strings.Join(args, "\n")
            if strings.TrimSpace(in) == "" {
                return fmt.Errorf("no input provided")
            }
        }
        lines := cliio.ReadLines([]byte(in))
        for _, line := range lines {
            t, err := dateutil.ParseFlexible(line, toEpochLayout, toEpochUTC)
            if err != nil {
                return err
            }
            if toEpochMs {
                fmt.Println(t.UnixMilli())
            } else {
                fmt.Println(t.Unix())
            }
        }
        return nil
    },
}

func init() {
    dateToEpochCmd.Flags().StringVar(&toEpochLayout, "layout", "", "Go time layout to parse (optional)")
    dateToEpochCmd.Flags().BoolVar(&toEpochUTC, "utc", false, "parse as UTC when timezone missing")
    dateToEpochCmd.Flags().BoolVar(&toEpochMs, "ms", false, "output milliseconds instead of seconds")
}
