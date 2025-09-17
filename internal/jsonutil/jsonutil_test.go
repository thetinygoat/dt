package jsonutil

import (
    "encoding/json"
    "strconv"
    "strings"
    "testing"
)

func TestMaybeUnquote(t *testing.T) {
    raw := `{"a":1}`
    s1 := strconv.Quote(raw)
    s2 := strconv.Quote(s1)
    s3 := strconv.Quote(s2)
    out := string(MaybeUnquote([]byte(s3)))
    if out != raw {
        t.Fatalf("expected fully unquoted to raw; got %q", out)
    }
}

func TestPretty_RawAndStringified(t *testing.T) {
    // Raw
    got, err := Pretty([]byte(`{"a":1}`), 2)
    if err != nil { t.Fatal(err) }
    if !strings.Contains(string(got), "\"a\": 1") {
        t.Fatalf("unexpected pretty output: %q", string(got))
    }
    // Stringified
    str := strconv.Quote(`{"a":1}`)
    got2, err := Pretty([]byte(str), 2)
    if err != nil { t.Fatal(err) }
    if !strings.Contains(string(got2), "\"a\": 1") {
        t.Fatalf("unexpected pretty output for stringified: %q", string(got2))
    }
}

func TestCompactMinify(t *testing.T) {
    in := []byte("{\n  \"a\": 1,\n  \"b\": [1, 2]\n}")
    got, err := CompactMinify(in, false)
    if err != nil { t.Fatal(err) }
    if string(got) != `{"a":1,"b":[1,2]}` {
        t.Fatalf("minify mismatch: %q", string(got))
    }
}

func TestStringify(t *testing.T) {
    in := []byte(`{"a":1}`)
    // Default (quoted)
    got, err := Stringify(in, false, false)
    if err != nil { t.Fatal(err) }
    if !strings.HasPrefix(string(got), "\"") || !strings.HasSuffix(string(got), "\"") {
        t.Fatalf("expected surrounding quotes: %q", string(got))
    }
    // The inner should be valid JSON
    var inner string
    if err := json.Unmarshal(got, &inner); err != nil {
        t.Fatalf("unmarshal string literal failed: %v", err)
    }
    if inner != `{"a":1}` { t.Fatalf("inner mismatch: %q", inner) }

    // No quotes
    got2, err := Stringify(in, false, true)
    if err != nil { t.Fatal(err) }
    if strings.HasPrefix(string(got2), "\"") || strings.HasSuffix(string(got2), "\"") {
        t.Fatalf("did not expect quotes: %q", string(got2))
    }
}
