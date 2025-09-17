package dateutil

import (
    "errors"
    "strconv"
    "strings"
    "time"
)

var CommonLayouts = []string{
    time.RFC3339Nano,
    time.RFC3339,
    time.RFC1123Z,
    time.RFC1123,
    time.RFC822Z,
    time.RFC822,
    time.ANSIC,
    time.UnixDate,
    time.RubyDate,
    "2006-01-02 15:04:05 MST",
    "2006-01-02 15:04:05 -0700",
    "2006-01-02 15:04:05",
    "2006-01-02 15:04 MST",
    "2006-01-02 15:04",
    "2006-01-02",
}

// ParseFlexible attempts multiple formats; if layout is provided, uses it first.
func ParseFlexible(s string, layout string, utc bool) (time.Time, error) {
    s = strings.TrimSpace(s)
    if s == "" {
        return time.Time{}, errors.New("empty input")
    }
    // numeric epoch detection as convenience
    if isAllDigits(s) {
        // detect ms vs s
        if len(s) > 10 {
            // milliseconds
            ms, _ := strconv.ParseInt(s, 10, 64)
            sec := ms / 1000
            nsec := (ms % 1000) * int64(time.Millisecond)
            return time.Unix(sec, nsec), nil
        }
        sec, _ := strconv.ParseInt(s, 10, 64)
        return time.Unix(sec, 0), nil
    }
    var loc *time.Location
    if utc {
        loc = time.UTC
    } else {
        loc = time.Local
    }
    // custom layout first
    if layout != "" {
        if t, err := time.ParseInLocation(layout, s, loc); err == nil {
            return t, nil
        }
    }
    // try common layouts
    for _, l := range CommonLayouts {
        if t, err := time.ParseInLocation(l, s, loc); err == nil {
            return t, nil
        }
    }
    return time.Time{}, errors.New("could not parse time; provide --layout")
}

func isAllDigits(s string) bool {
    for _, r := range s {
        if r < '0' || r > '9' {
            return false
        }
    }
    return true
}

// FormatTime renders t in selected format.
// format: "rfc3339", "unix", "unixms", or "layout" (with layout value).
func FormatTime(t time.Time, format string, layout string, utc bool) string {
    if utc {
        t = t.UTC()
    }
    switch strings.ToLower(format) {
    case "", "rfc3339":
        return t.Format(time.RFC3339)
    case "unix":
        return strconv.FormatInt(t.Unix(), 10)
    case "unixms":
        return strconv.FormatInt(t.UnixMilli(), 10)
    case "layout":
        if layout == "" {
            layout = time.RFC3339
        }
        return t.Format(layout)
    default:
        return t.Format(format)
    }
}

