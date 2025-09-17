package cmd

import (
    "encoding/base64"
    "fmt"
    "strings"

    "dt/internal/cliio"
    "github.com/spf13/cobra"
)

var (
    b64url   bool
    b64nopad bool
)

func init() {
    base64Cmd.AddCommand(base64EncodeCmd)
    base64Cmd.AddCommand(base64DecodeCmd)
    rootCmd.AddCommand(base64Cmd)
}

var base64Cmd = &cobra.Command{Use: "base64", Short: "Base64 encode/decode"}

var base64EncodeCmd = &cobra.Command{
    Use:   "encode",
    Short: "Base64-encode input",
    RunE: func(cmd *cobra.Command, args []string) error {
        in, err := cliio.ReadAll(args)
        if err != nil {
            return err
        }
        var enc *base64.Encoding
        if b64url {
            if b64nopad {
                enc = base64.RawURLEncoding
            } else {
                enc = base64.URLEncoding
            }
        } else {
            if b64nopad {
                enc = base64.RawStdEncoding
            } else {
                enc = base64.StdEncoding
            }
        }
        out := enc.EncodeToString(in)
        fmt.Println(out)
        return nil
    },
}

var base64DecodeCmd = &cobra.Command{
    Use:   "decode",
    Short: "Base64-decode input",
    RunE: func(cmd *cobra.Command, args []string) error {
        in, err := cliio.ReadAll(args)
        if err != nil {
            return err
        }
        s := strings.TrimSpace(string(in))
        // Try both raw and padded variants; URL if --url set.
        var dec *base64.Encoding
        try := func(e *base64.Encoding) ([]byte, error) { return e.DecodeString(s) }
        if b64url {
            dec = base64.RawURLEncoding
            if out, err := try(dec); err == nil {
                fmt.Println(string(out))
                return nil
            }
            dec = base64.URLEncoding
            if out, err := try(dec); err == nil {
                fmt.Println(string(out))
                return nil
            }
        } else {
            dec = base64.RawStdEncoding
            if out, err := try(dec); err == nil {
                fmt.Println(string(out))
                return nil
            }
            dec = base64.StdEncoding
            if out, err := try(dec); err == nil {
                fmt.Println(string(out))
                return nil
            }
        }
        return fmt.Errorf("invalid base64 input")
    },
}

func init() {
    base64EncodeCmd.Flags().BoolVar(&b64url, "url", false, "use URL-safe encoding")
    base64EncodeCmd.Flags().BoolVar(&b64nopad, "no-pad", false, "omit '=' padding")
    base64DecodeCmd.Flags().BoolVar(&b64url, "url", false, "expect URL-safe encoding variants")
}

