package cmd

import (
    "fmt"
    "strings"
    "time"

    "dt/internal/dateutil"
    "github.com/spf13/cobra"
)

var (
    addDuration string
    addFrom     string
    addUTC      bool
    addFormat   string
    addLayout   string
)

func init() {
    dateCmd.AddCommand(dateAddCmd)
}

var dateAddCmd = &cobra.Command{
    Use:   "add",
    Short: "Add a duration to now or a given time",
    Long:  "Adds a Go-style duration (e.g., 90m, 1h30m, 48h) to now or a provided --from value (parsed flexibly).",
    RunE: func(cmd *cobra.Command, args []string) error {
        if strings.TrimSpace(addDuration) == "" {
            return fmt.Errorf("--duration is required")
        }
        dur, err := time.ParseDuration(addDuration)
        if err != nil {
            return err
        }
        var base time.Time
        if strings.TrimSpace(addFrom) == "" {
            base = time.Now()
        } else {
            base, err = dateutil.ParseFlexible(addFrom, "", addUTC)
            if err != nil {
                return err
            }
        }
        t := base.Add(dur)
        out := dateutil.FormatTime(t, addFormat, addLayout, addUTC)
        fmt.Println(out)
        return nil
    },
}

func init() {
    dateAddCmd.Flags().StringVar(&addDuration, "duration", "", "Go duration to add, e.g., 1h30m")
    dateAddCmd.Flags().StringVar(&addFrom, "from", "", "optional base time or epoch")
    dateAddCmd.Flags().BoolVar(&addUTC, "utc", false, "treat base/print as UTC")
    dateAddCmd.Flags().StringVar(&addFormat, "format", "rfc3339", "output format: rfc3339|unix|unixms|layout|<Go layout>")
    dateAddCmd.Flags().StringVar(&addLayout, "layout", "", "when --format=layout, Go time layout")
}
