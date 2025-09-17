package cmd

import (
    "fmt"
    "strconv"
    "strings"
    "time"

    "dt/internal/cliio"
    "dt/internal/dateutil"
    "github.com/spf13/cobra"
)

var (
    fromEpochFormat string
    fromEpochLayout string
    fromEpochUTC    bool
)

func init() {
    dateCmd.AddCommand(dateFromEpochCmd)
}

var dateFromEpochCmd = &cobra.Command{
    Use:   "from-epoch [values...]",
    Short: "Convert epoch to human time",
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
            s := strings.TrimSpace(line)
            if s == "" {
                continue
            }
            // Detect seconds vs milliseconds
            n, err := strconv.ParseInt(s, 10, 64)
            if err != nil {
                return fmt.Errorf("invalid epoch value: %s", s)
            }
            var t time.Time
            if len(s) > 10 {
                t = time.Unix(0, n*int64(time.Millisecond))
            } else {
                t = time.Unix(n, 0)
            }
            out := dateutil.FormatTime(t, fromEpochFormat, fromEpochLayout, fromEpochUTC)
            fmt.Println(out)
        }
        return nil
    },
}

func init() {
    dateFromEpochCmd.Flags().StringVar(&fromEpochFormat, "format", "rfc3339", "output format: rfc3339|unix|unixms|layout|<Go layout>")
    dateFromEpochCmd.Flags().StringVar(&fromEpochLayout, "layout", "", "Go time layout when --format=layout")
    dateFromEpochCmd.Flags().BoolVar(&fromEpochUTC, "utc", false, "print in UTC")
    dateFromEpochCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
        return []string{"rfc3339", "unix", "unixms", "layout"}, cobra.ShellCompDirectiveNoFileComp
    })
}

