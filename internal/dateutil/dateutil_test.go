package dateutil

import (
    "strconv"
    "testing"
    "time"
)

func TestParseFlexibleAndFormat(t *testing.T) {
    // Parse a naive datetime; expect same string when formatted with the same layout.
    s := "2025-09-17 12:34:56"
    got, err := ParseFlexible(s, "", false)
    if err != nil { t.Fatal(err) }
    if FormatTime(got, "layout", "2006-01-02 15:04:05", false) != s {
        t.Fatalf("round-trip mismatch")
    }

    // Parse RFC3339 Zulu and format to unix
    s2 := "1970-01-01T00:00:00Z"
    t2, err := ParseFlexible(s2, "", true)
    if err != nil { t.Fatal(err) }
    if FormatTime(t2, "unix", "", true) != "0" {
        t.Fatalf("expected 0 epoch")
    }
}

func TestEpochDetection(t *testing.T) {
    // Seconds
    t1, err := ParseFlexible("1", "", true)
    if err != nil { t.Fatal(err) }
    if t1.Unix() != 1 { t.Fatalf("expected 1 second") }
    // Milliseconds (13-digit)
    const msStr = "1690000000123"
    t2, err := ParseFlexible(msStr, "", true)
    if err != nil { t.Fatal(err) }
    if got := strconv.FormatInt(t2.UnixMilli(), 10); got != msStr { t.Fatalf("expected %s, got %s", msStr, got) }
}

func TestFormatVariants(t *testing.T) {
    ref := time.Unix(42, 0).UTC()
    if FormatTime(ref, "unix", "", true) != "42" { t.Fatalf("unix mismatch") }
    if FormatTime(ref, "unixms", "", true) != strconv.FormatInt(ref.UnixMilli(), 10) { t.Fatalf("unixms mismatch") }
    if FormatTime(ref, "rfc3339", "", true) != "1970-01-01T00:00:42Z" { t.Fatalf("rfc3339 mismatch") }
}
