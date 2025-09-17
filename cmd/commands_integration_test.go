package cmd

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// run executes the Cobra root command with given args and optional stdin data.
// It captures stdout/stderr for assertions.
func run(t *testing.T, args []string, stdin string) (string, string, error) {
	t.Helper()
	// Capture stdout/stderr
	oldOut, oldErr, oldIn := os.Stdout, os.Stderr, os.Stdin
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout, os.Stderr = wOut, wErr
	// Prepare stdin (pipe to look like real piping)
	if stdin != "" {
		rIn, wIn, _ := os.Pipe()
		os.Stdin = rIn
		go func() {
			io.WriteString(wIn, stdin)
			wIn.Close()
		}()
	}
	// Ensure Cobra doesn’t print usage on errors when we assert them
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	// Restore
	wOut.Close()
	wErr.Close()
	os.Stdout, os.Stderr, os.Stdin = oldOut, oldErr, oldIn
	var bufOut, bufErr bytes.Buffer
	io.Copy(&bufOut, rOut)
	io.Copy(&bufErr, rErr)
	return bufOut.String(), bufErr.String(), err
}

func TestJSONPretty_WithPipedStringified(t *testing.T) {
	input := strconv.Quote(`{"a":1,"b":[1,2]}`)
	out, _, err := run(t, []string{"json", "pretty"}, input)
	if err != nil {
		t.Fatalf("cmd error: %v", err)
	}
	if !strings.Contains(out, "\"a\": 1") || !strings.Contains(out, "\n  \"b\": [\n    1,\n    2\n  ]") {
		t.Fatalf("unexpected pretty: %q", out)
	}
}

func TestJSONStringify_DefaultAndNoQuotes(t *testing.T) {
	out, _, err := run(t, []string{"json", "stringify", `{"a":1}`}, "")
	if err != nil {
		t.Fatalf("cmd error: %v", err)
	}
	if !strings.HasPrefix(out, "\"") || !strings.Contains(out, `\"a\":1`) {
		t.Fatalf("unexpected stringify: %q", out)
	}
	out2, _, err := run(t, []string{"json", "stringify", `{"a":1}`, "--no-quotes"}, "")
	if err != nil {
		t.Fatalf("cmd error: %v", err)
	}
	if strings.HasPrefix(strings.TrimSpace(out2), "\"") {
		t.Fatalf("expected no quotes: %q", out2)
	}
}

func TestBase64_EncodeDecode(t *testing.T) {
	enc, _, err := run(t, []string{"base64", "encode"}, "hello")
	if err != nil {
		t.Fatalf("encode err: %v", err)
	}
	if strings.TrimSpace(enc) != "aGVsbG8=" {
		t.Fatalf("unexpected b64: %q", enc)
	}
	dec, _, err := run(t, []string{"base64", "decode"}, strings.TrimSpace(enc))
	if err != nil {
		t.Fatalf("decode err: %v", err)
	}
	if strings.TrimSpace(dec) != "hello" {
		t.Fatalf("unexpected decoded: %q", dec)
	}
}

func TestDate_Conversions(t *testing.T) {
	// to-epoch
	to, _, err := run(t, []string{"date", "to-epoch", "--utc"}, "1970-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("to-epoch err: %v", err)
	}
	if strings.TrimSpace(to) != "0" {
		t.Fatalf("expected 0: %q", to)
	}
	// from-epoch
	from, _, err := run(t, []string{"date", "from-epoch", "--format", "rfc3339", "--utc"}, "0")
	if err != nil {
		t.Fatalf("from-epoch err: %v", err)
	}
	if strings.TrimSpace(from) != "1970-01-01T00:00:00Z" {
		t.Fatalf("unexpected: %q", from)
	}
	// add 1500ms from epoch
	add, _, err := run(t, []string{"date", "add", "--duration", "1500ms", "--from", "1970-01-01T00:00:00Z", "--format", "unixms", "--utc"}, "")
	if err != nil {
		t.Fatalf("add err: %v", err)
	}
	if strings.TrimSpace(add) != "1500" {
		t.Fatalf("unexpected add: %q", add)
	}
}

func TestUUID_New(t *testing.T) {
	out, _, err := run(t, []string{"uuid", "new", "-n", "3"}, "")
	if err != nil {
		t.Fatalf("uuid err: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 UUIDs: %v", lines)
	}
	re := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	for _, l := range lines {
		if !re.MatchString(l) {
			t.Fatalf("invalid uuid: %q", l)
		}
	}
}

func TestEnv_FromJSON_Flatten(t *testing.T) {
	in := `{"db":{"name":"x"},"port":8080}`
	out, _, err := run(t, []string{"env", "from-json", "--uppercase", "--flatten", "--sep", "_", "--prefix", "APP_"}, in)
	if err != nil {
		t.Fatalf("env err: %v", err)
	}
	// Order is sorted; assert expected lines are present
	want := []string{"APP_DB_NAME=x", "APP_PORT=8080"}
	for _, w := range want {
		if !strings.Contains(out, w+"\n") {
			t.Fatalf("missing %q in %q", w, out)
		}
	}
}

func TestEnv_FromKV(t *testing.T) {
	in := "host: localhost\nport: 8080\n# c\n"
	out, _, err := run(t, []string{"env", "from-kv", "--uppercase", "--prefix", "APP_"}, in)
	if err != nil {
		t.Fatalf("env kv err: %v", err)
	}
	if !strings.Contains(out, "APP_HOST=localhost\n") || !strings.Contains(out, "APP_PORT=8080\n") {
		t.Fatalf("unexpected kv output: %q", out)
	}
}

func TestTextJoin_DefaultSingleQuote(t *testing.T) {
	in := "Alice\nBob\n"
	out, _, err := run(t, []string{"text", "join"}, in)
	if err != nil {
		t.Fatalf("text join err: %v", err)
	}
	if strings.TrimSpace(out) != "'Alice','Bob'" {
		t.Fatalf("unexpected joined output: %q", out)
	}
}

func TestTextJoin_DoubleQuoteCSV(t *testing.T) {
	in := "Alice,Bob,\"Charlie, Q.\"\n"
	out, _, err := run(t, []string{"text", "join", "--quote", "double", "--split", "csv"}, in)
	if err != nil {
		t.Fatalf("text join csv err: %v", err)
	}
	if strings.TrimSpace(out) != "\"Alice\",\"Bob\",\"Charlie, Q.\"" {
		t.Fatalf("unexpected csv join: %q", out)
	}
}

func TestTextJoin_CustomSepUnique(t *testing.T) {
	in := "Alpha\nBeta\nAlpha\n"
	out, _, err := run(t, []string{"text", "join", "--sep", "|", "--unique", "--quote", "single"}, in)
	if err != nil {
		t.Fatalf("text join unique err: %v", err)
	}
	if strings.TrimSpace(out) != "'Alpha'|'Beta'" {
		t.Fatalf("unexpected unique join: %q", out)
	}
}

func TestHash_Digests(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		stdin string
		want  string
	}{
		{
			name:  "sha256 hex",
			args:  []string{"hash", "sha256"},
			stdin: "hello",
			want:  "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:  "sha256 base64",
			args:  []string{"hash", "sha256", "--encoding", "base64"},
			stdin: "hello",
			want:  "LPJNul+wow4m6DsqxbninhsWHlwfp0JecwQzYpOLmCQ=",
		},
		{
			name:  "md5 salted",
			args:  []string{"hash", "md5", "--salt", "pepper"},
			stdin: "hello",
			want:  "6967321c83e9f01a33e7edecce748877",
		},
		{
			name:  "sha3-256 hex",
			args:  []string{"hash", "sha3-256"},
			stdin: "hello",
			want:  "3338be694f50c5f338814986cdf0686453a888b84f424d792af4b9202398f392",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, _, err := run(t, tc.args, tc.stdin)
			if err != nil {
				t.Fatalf("hash err: %v", err)
			}
			if got := strings.TrimSpace(out); got != tc.want {
				t.Fatalf("unexpected digest: got %q want %q", got, tc.want)
			}
		})
	}
}
