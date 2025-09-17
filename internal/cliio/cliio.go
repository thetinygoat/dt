package cliio

import (
    "bufio"
    "bytes"
    "errors"
    "io"
    "os"
    "strings"
)

// IsInputFromPipe reports whether stdin is a named pipe / not a TTY.
func IsInputFromPipe() bool {
    fi, err := os.Stdin.Stat()
    if err != nil {
        return false
    }
    return (fi.Mode() & os.ModeCharDevice) == 0
}

// ReadAll reads from stdin if piped, otherwise joins args with spaces.
func ReadAll(args []string) ([]byte, error) {
    if IsInputFromPipe() {
        return io.ReadAll(bufio.NewReader(os.Stdin))
    }
    if len(args) == 0 {
        return nil, errors.New("no input provided; pass arguments or pipe data")
    }
    return []byte(strings.Join(args, " ")), nil
}

// ReadLines splits input into lines, trimming trailing CRLF.
func ReadLines(b []byte) []string {
    s := string(bytes.ReplaceAll(bytes.TrimSpace(b), []byte("\r"), nil))
    if s == "" {
        return nil
    }
    return strings.Split(s, "\n")
}

// Println writes s + newline to stdout.
func Println(s string) error {
    _, err := os.Stdout.WriteString(s + "\n")
    return err
}

