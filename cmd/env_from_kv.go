package cmd

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "sort"
    "strings"

    "dt/internal/cliio"
    "github.com/spf13/cobra"
)

var (
    kvUpper  bool
    kvPrefix string
)

func init() {
    envFromKVCmd.Flags().BoolVar(&kvUpper, "uppercase", false, "upper-case keys")
    envFromKVCmd.Flags().StringVar(&kvPrefix, "prefix", "", "prefix to add to each key")
}

var envFromKVCmd = &cobra.Command{
    Use:   "from-kv",
    Short: "Convert key:value lines to KEY=VALUE lines",
    Long:  "Reads lines like 'key: value' (JSON/YAML style) and prints 'KEY=VALUE'. Ignores blank lines and '#' comments.",
    RunE: func(cmd *cobra.Command, args []string) error {
        var r *bufio.Reader
        if cliio.IsInputFromPipe() {
            r = bufio.NewReader(os.Stdin)
        } else {
            // treat args as newline-separated input
            s := strings.Join(args, "\n")
            r = bufio.NewReader(strings.NewReader(s))
        }
        type pair struct{ k, v string }
        pairs := []pair{}
        re := regexp.MustCompile(`^\s*([^:#\s][^:]*)\s*:\s*(.*)$`)
        for {
            line, err := r.ReadString('\n')
            if len(line) > 0 {
                l := strings.TrimSpace(line)
                if l == "" || strings.HasPrefix(l, "#") { continue }
                if m := re.FindStringSubmatch(l); m != nil {
                    k := strings.TrimSpace(m[1])
                    v := strings.TrimSpace(m[2])
                    if kvUpper { k = strings.ToUpper(k) }
                    if kvPrefix != "" { k = kvPrefix + k }
                    pairs = append(pairs, pair{k, v})
                } else {
                    return fmt.Errorf("invalid line: %q (expected 'key: value')", l)
                }
            }
            if err != nil {
                break
            }
        }
        sort.Slice(pairs, func(i, j int) bool { return pairs[i].k < pairs[j].k })
        for _, p := range pairs {
            fmt.Printf("%s=%s\n", p.k, p.v)
        }
        return nil
    },
}

