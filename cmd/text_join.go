package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	"dt/internal/cliio"
	"github.com/spf13/cobra"
)

var (
	textJoinSep       string
	textJoinQuote     string
	textJoinSplit     string
	textJoinTrim      bool
	textJoinSkipEmpty bool
	textJoinUnique    bool
)

func init() {
	rootCmd.AddCommand(textCmd)
	textCmd.AddCommand(textJoinCmd)

	textJoinCmd.Flags().StringVar(&textJoinSep, "sep", ",", "separator between items (supports \\n, \\t, \\r, \\0 escapes)")
	textJoinCmd.Flags().StringVar(&textJoinQuote, "quote", "single", "quote style: single|double|none")
	textJoinCmd.Flags().StringVar(&textJoinSplit, "split", "lines", "input splitter: lines|tab|csv")
	textJoinCmd.Flags().BoolVar(&textJoinTrim, "trim", true, "trim whitespace around each item")
	textJoinCmd.Flags().BoolVar(&textJoinSkipEmpty, "skip-empty", true, "drop empty items after trimming")
	textJoinCmd.Flags().BoolVar(&textJoinUnique, "unique", false, "deduplicate items (first occurrence wins)")

	textJoinCmd.RegisterFlagCompletionFunc("quote", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"single", "double", "none"}, cobra.ShellCompDirectiveNoFileComp
	})
	textJoinCmd.RegisterFlagCompletionFunc("split", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"lines", "tab", "csv"}, cobra.ShellCompDirectiveNoFileComp
	})
}

var textCmd = &cobra.Command{Use: "text", Short: "Text utilities"}

var textJoinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join rows or columns into a single separator-delimited line",
	Long: `Reads newline/tab/CSV separated input and emits one line joined by a custom separator.
Reads stdin when piped; otherwise treats positional arguments as individual items.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		items, err := collectTextItems(args)
		if err != nil {
			return err
		}
		if len(items) == 0 {
			return errors.New("no input provided")
		}
		sep, err := interpretEscapes(textJoinSep)
		if err != nil {
			return err
		}
		formatted, err := formatItems(items)
		if err != nil {
			return err
		}
		return cliio.Println(strings.Join(formatted, sep))
	},
}

func collectTextItems(args []string) ([]string, error) {
	var raw string
	if cliio.IsInputFromPipe() {
		b, err := cliio.ReadAll(nil)
		if err != nil {
			return nil, err
		}
		raw = string(b)
	} else if len(args) > 0 {
		raw = strings.Join(args, "\n")
	} else {
		return nil, errors.New("no input provided")
	}

	switch textJoinSplit {
	case "lines":
		return splitLines(raw), nil
	case "tab":
		return splitTab(raw), nil
	case "csv":
		return splitCSV(raw)
	default:
		return nil, fmt.Errorf("unsupported split mode: %s", textJoinSplit)
	}
}

func splitLines(raw string) []string {
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")
	if raw == "" {
		return nil
	}
	return strings.Split(raw, "\n")
}

func splitTab(raw string) []string {
	if raw == "" {
		return nil
	}
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return r == '\n' || r == '\t'
	})
	return fields
}

func splitCSV(raw string) ([]string, error) {
	rdr := csv.NewReader(strings.NewReader(raw))
	rdr.FieldsPerRecord = -1
	var out []string
	for {
		rec, err := rdr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		out = append(out, rec...)
	}
	return out, nil
}

func interpretEscapes(s string) (string, error) {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' {
			b.WriteByte(s[i])
			continue
		}
		if i+1 >= len(s) {
			return "", errors.New("trailing backslash in separator")
		}
		i++
		switch s[i] {
		case 'n':
			b.WriteByte('\n')
		case 't':
			b.WriteByte('\t')
		case 'r':
			b.WriteByte('\r')
		case '0':
			b.WriteByte(0)
		case '\\':
			b.WriteByte('\\')
		default:
			b.WriteByte('\\')
			b.WriteByte(s[i])
		}
	}
	return b.String(), nil
}

func formatItems(items []string) ([]string, error) {
	var out []string
	seen := map[string]struct{}{}
	for _, item := range items {
		if textJoinTrim {
			item = strings.TrimSpace(item)
		}
		if textJoinSkipEmpty && item == "" {
			continue
		}
		key := item
		if textJoinUnique {
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
		}
		q, err := quoteItem(item)
		if err != nil {
			return nil, err
		}
		out = append(out, q)
	}
	return out, nil
}

func quoteItem(s string) (string, error) {
	switch textJoinQuote {
	case "single":
		return "'" + strings.ReplaceAll(s, "'", "''") + "'", nil
	case "double":
		return "\"" + strings.ReplaceAll(s, "\"", "\\\"") + "\"", nil
	case "none":
		return s, nil
	default:
		return "", fmt.Errorf("invalid quote style: %s", textJoinQuote)
	}
}
