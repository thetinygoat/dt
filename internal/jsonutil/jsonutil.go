package jsonutil

import (
    "bytes"
    "encoding/json"
    "errors"
    "strconv"
    "strings"
)

// MaybeUnquote tries to unquote input up to 3 times if it's a quoted string.
func MaybeUnquote(in []byte) []byte {
    s := strings.TrimSpace(string(in))
    for i := 0; i < 3; i++ {
        if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
            u, err := strconv.Unquote(s)
            if err != nil {
                break
            }
            s = u
        } else {
            break
        }
    }
    return []byte(s)
}

// Pretty formats JSON data with indentation.
func Pretty(in []byte, indent int) ([]byte, error) {
    in = bytes.TrimSpace(in)
    if len(in) == 0 {
        return nil, errors.New("empty input")
    }
    sp := strings.Repeat(" ", indent)
    var out bytes.Buffer

    // If input starts with a quote, it may be a JSON string literal containing JSON.
    if len(in) > 0 && in[0] == '"' {
        uq := MaybeUnquote(in)
        if len(uq) > 0 && (uq[0] == '{' || uq[0] == '[') {
            if err := json.Indent(&out, uq, "", sp); err == nil {
                return out.Bytes(), nil
            }
        }
        // Fallthrough to try raw indent as a last resort.
    }

    // Try raw JSON indent for objects/arrays.
    if len(in) > 0 && (in[0] == '{' || in[0] == '[') {
        if err := json.Indent(&out, in, "", sp); err == nil {
            return out.Bytes(), nil
        }
    }

    // Try unquoted anyway (handles nested quotes cases).
    uq := MaybeUnquote(in)
    if len(uq) > 0 && (uq[0] == '{' || uq[0] == '[') {
        out.Reset()
        if err := json.Indent(&out, uq, "", sp); err == nil {
            return out.Bytes(), nil
        }
    }
    return nil, errors.New("invalid JSON or stringified JSON")
}

// CompactMinify minifies JSON input; optionally unquotes first.
func CompactMinify(in []byte, allowUnquote bool) ([]byte, error) {
    in = bytes.TrimSpace(in)
    if len(in) == 0 {
        return nil, errors.New("empty input")
    }
    var out bytes.Buffer
    if json.Compact(&out, in) == nil {
        return out.Bytes(), nil
    }
    if allowUnquote {
        out.Reset()
        uq := MaybeUnquote(in)
        if json.Compact(&out, uq) == nil {
            return out.Bytes(), nil
        }
    }
    return nil, errors.New("invalid JSON")
}

// Stringify returns a JSON string literal of the given JSON input.
// If compact is true, input is first minified; if noQuotes is true, surrounding quotes are stripped.
func Stringify(in []byte, compact bool, noQuotes bool) ([]byte, error) {
    var data []byte
    var err error
    if compact {
        data, err = CompactMinify(in, true)
    } else {
        // Attempt raw or unquoted parse to ensure it's valid JSON; then re-marshal
        var v any
        if e := json.Unmarshal(in, &v); e != nil {
            uq := MaybeUnquote(in)
            if e2 := json.Unmarshal(uq, &v); e2 != nil {
                return nil, errors.New("invalid JSON")
            }
        }
        data, err = json.Marshal(v) // canonical compact form
    }
    if err != nil {
        return nil, err
    }
    // Encode as string literal
    s := strconv.Quote(string(data))
    if noQuotes {
        // Remove surrounding quotes
        return []byte(s[1 : len(s)-1]), nil
    }
    return []byte(s), nil
}
