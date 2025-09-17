package cmd

import (
    "crypto/rand"
    "fmt"

    "github.com/spf13/cobra"
)

var uuidCount int

func init() {
    rootCmd.AddCommand(uuidCmd)
    uuidCmd.AddCommand(uuidNewCmd)
}

var uuidCmd = &cobra.Command{Use: "uuid", Short: "UUID utilities"}

var uuidNewCmd = &cobra.Command{
    Use:   "new",
    Short: "Generate UUIDv4",
    RunE: func(cmd *cobra.Command, args []string) error {
        if uuidCount <= 0 {
            uuidCount = 1
        }
        for i := 0; i < uuidCount; i++ {
            b := make([]byte, 16)
            if _, err := rand.Read(b); err != nil {
                return err
            }
            // Set version (4) and variant (RFC 4122)
            b[6] = (b[6] & 0x0f) | 0x40
            b[8] = (b[8] & 0x3f) | 0x80
            fmt.Printf("%x-%x-%x-%x-%x\n", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
        }
        return nil
    },
}

func init() {
    uuidNewCmd.Flags().IntVarP(&uuidCount, "count", "n", 1, "number of UUIDs to generate")
}

