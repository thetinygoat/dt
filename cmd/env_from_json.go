package cmd

import (
    "encoding/json"
    "fmt"
    "sort"
    "strings"

    "dt/internal/cliio"
    "github.com/spf13/cobra"
)

var (
    envUpper    bool
    envPrefix   string
    envFlatten  bool
    envSep      string
)

func init() {
    rootCmd.AddCommand(envCmd)
    envCmd.AddCommand(envFromJSONCmd)
    envCmd.AddCommand(envFromKVCmd)
}

var envCmd = &cobra.Command{Use: "env", Short: "Environment helpers"}

var envFromJSONCmd = &cobra.Command{
    Use:   "from-json",
    Short: "Convert a JSON object to KEY=VALUE lines",
    RunE: func(cmd *cobra.Command, args []string) error {
        b, err := cliio.ReadAll(args)
        if err != nil {
            return err
        }
        var v any
        if err := json.Unmarshal(b, &v); err != nil {
            return fmt.Errorf("expected a JSON object: %w", err)
        }
        obj, ok := v.(map[string]any)
        if !ok {
            return fmt.Errorf("expected a JSON object at top-level")
        }
        pairs := make(map[string]string)
        if envFlatten {
            flattenJSON(obj, "", pairs)
        } else {
            for k, vv := range obj {
                pairs[normalizeKey(k)] = stringifySimple(vv)
            }
        }
        // order keys for stable output
        keys := make([]string, 0, len(pairs))
        for k := range pairs {
            keys = append(keys, k)
        }
        sort.Strings(keys)
        for _, k := range keys {
            key := k
            if envUpper { key = strings.ToUpper(key) }
            if envPrefix != "" { key = envPrefix + key }
            fmt.Printf("%s=%s\n", key, pairs[k])
        }
        return nil
    },
}

func init() {
    envFromJSONCmd.Flags().BoolVar(&envUpper, "uppercase", false, "upper-case keys")
    envFromJSONCmd.Flags().StringVar(&envPrefix, "prefix", "", "prefix to add to each key")
    envFromJSONCmd.Flags().BoolVar(&envFlatten, "flatten", false, "flatten nested objects")
    envFromJSONCmd.Flags().StringVar(&envSep, "sep", "_", "separator for flattened keys")
}

func normalizeKey(k string) string {
    // Replace spaces and dots with sep when flattening later; basic cleanup
    return strings.TrimSpace(k)
}

func stringifySimple(v any) string {
    switch t := v.(type) {
    case nil:
        return ""
    case string:
        return t
    case bool:
        if t { return "true" }
        return "false"
    case float64:
        // Preserve integer if possible
        if t == float64(int64(t)) {
            return fmt.Sprintf("%d", int64(t))
        }
        return fmt.Sprintf("%v", t)
    default:
        // Marshal nested structures as compact JSON string
        b, _ := json.Marshal(t)
        return string(b)
    }
}

func flattenJSON(m map[string]any, prefix string, out map[string]string) {
    for k, v := range m {
        key := normalizeKey(k)
        if prefix != "" { key = prefix + envSep + key }
        switch t := v.(type) {
        case map[string]any:
            flattenJSON(t, key, out)
        case []any:
            // arrays -> JSON string to preserve order
            out[key] = stringifySimple(t)
        default:
            out[key] = stringifySimple(t)
        }
    }
}

