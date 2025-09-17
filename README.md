# dt — Day-to-day Dev Toolbox

A small, fast Go CLI that bundles everyday developer utilities—JSON inspection, date conversions, UUID generation, Base64 transforms, and environment-variable helpers—in a single command built for shell piping.

---

## Install

```sh
go install ./...
# or
go build -o dt .
```

_Requires Go 1.21 or newer. The first invocation pulls `github.com/spf13/cobra` automatically._

---

## Usage Basics

- Every command reads from stdin when piping is detected; otherwise it uses positional arguments.
- Successful results go to stdout; diagnostics go to stderr so you can chain commands safely.
- Run `dt --help`, `dt <namespace> --help`, or `dt <namespace> <command> --help` for interactive flag info.

### Quick Glance

```sh
# Pretty-print JSON (stringified or raw)
echo '"{\"a\":1,\"b\":[1,2]}"' | dt json pretty

# Convert JSON to a single escaped string literal
dt json stringify '{"project":"dt","version":1}'

# Round-trip Base64
echo -n 'hello' | dt base64 encode | dt base64 decode

# Epoch and layout conversions
dt date to-epoch '2025-09-17T12:34:56Z'
dt date from-epoch --format layout --layout '2006-01-02 15:04:05' 1758102896

dt uuid new -n 3

# Environment helpers
cat cfg.json | dt env from-json --uppercase --flatten --sep '_' --prefix APP_
printf 'host: localhost\nport: 8080\n' | dt env from-kv

# Join spreadsheet columns into a quoted list
pbpaste | dt text join --quote double --sep ', '
```

---

## Command Reference

All examples assume the binary name is `dt` and that you are running them from a Unix-like shell. Replace sample inputs with your own values.

### JSON Commands

#### `dt json pretty`

- **Synopsis:** `dt json pretty [--indent 2]`
- **Purpose:** Reformat JSON with configurable indentation. Automatically unwraps up to three levels of quoted JSON so you can pipe stringified payloads directly.
- **Flags:**
  - `--indent <n>` — number of spaces per level (default: `2`).
- **Example:**
  ```sh
  echo '"{\"service\":\"api\",\"ports\":[80,443]}"' | dt json pretty
  # Output
  # {
  #   "service": "api",
  #   "ports": [
  #     80,
  #     443
  #   ]
  # }
  ```

#### `dt json stringify`

- **Synopsis:** `dt json stringify [--compact] [--no-quotes] <json|stdin>`
- **Purpose:** Turn JSON into a JSON string literal that can be embedded in config or code.
- **Flags:**
  - `--compact` — emit the shortest equivalent string (no spaces or newlines).
  - `--no-quotes` — omit the surrounding quotes (useful when another tool adds its own quoting).
- **Example:**

  ```sh
  dt json stringify '{"scope":"deploy","features":["json","date"]}'
  # Output
  # "{\"scope\":\"deploy\",\"features\":[\"json\",\"date\"]}"

  echo '{"scope":"deploy"}' | dt json stringify --compact --no-quotes
  # Output
  # {\"scope\":\"deploy\"}
  ```

### Base64 Commands

#### `dt base64 encode`

- **Synopsis:** `dt base64 encode [--url] [--no-pad]`
- **Purpose:** Convert stdin or arguments into Base64. Supports standard and URL-safe alphabets.
- **Flags:**
  - `--url` — use the URL-safe alphabet (`-` and `_`).
  - `--no-pad` — omit `=` padding (useful for JWT-style payloads).
- **Example:**

  ```sh
  echo -n 'hello' | dt base64 encode
  # Output
  # aGVsbG8=

  echo -n 'payload' | dt base64 encode --url --no-pad
  # Output
  # cGF5bG9hZA
  ```

#### `dt base64 decode`

- **Synopsis:** `dt base64 decode [--url]`
- **Purpose:** Decode Base64 strings, automatically trying padded and unpadded forms.
- **Flags:**
  - `--url` — treat input as URL-safe Base64.
- **Example:**
  ```sh
  echo 'YXBwOnNlY3JldA==' | dt base64 decode
  # Output
  # app:secret
  ```

### Text Commands

#### `dt text join`

- **Synopsis:** `dt text join [--sep ","] [--quote single|double|none] [--split lines|tab|csv] [--trim] [--skip-empty] [--unique] [items...]`
- **Purpose:** Collapse multi-line, tabular, or CSV input into a single separator-delimited row with optional quoting. Ideal for turning spreadsheet columns or clipboard lists into shell-ready or SQL-ready strings.
- **Flags:**
  - `--sep` — separator string; supports escape sequences such as `\n`, `\t`, `\r`, and NUL (`\0`). Default `,`.
  - `--quote` — wrap each value with single quotes (default), double quotes, or no quoting.
  - `--split` — choose how to split raw input: newline-delimited (`lines`), tab/line separated (`tab`), or CSV-aware (`csv`, respecting quoted commas).
  - `--trim` / `--no-trim` — control per-item whitespace trimming (defaults to on).
  - `--skip-empty` / `--skip-empty=false` — drop empty items after trimming (defaults to on).
  - `--unique` — keep the first occurrence of each value and drop duplicates.
- **Examples:**

  ```sh
  # Google Sheets column -> SQL IN clause
  pbpaste | dt text join --quote double --sep ', '
  # "Alice", "Bob", "Charlie"

  # TSV clipboard -> pipe-delimited list without quotes
  printf 'alpha\tbeta\tgamma\n' | dt text join --split tab --quote none --sep ' | '
  # alpha | beta | gamma

  # CSV row with embedded commas -> single-quoted values
  printf '"Widget, Large",Small\n' | dt text join --split csv
  # 'Widget, Large','Small'
  ```


### Date Commands

All date commands accept either CLI arguments or piped input. Layout strings use Go's `time` reference layout (`2006-01-02 15:04:05`).

#### `dt date now`

- **Synopsis:** `dt date now [--format rfc3339|unix|unixms|layout] [--layout <fmt>] [--utc]`
- **Purpose:** Print the current time. Defaults to RFC3339 in local time.
- **Example:**

  ```sh
  dt date now --utc
  # Output (example)
  # 2025-09-17T18:42:10Z

  dt date now --format unixms
  # Output (example)
  # 1758134530123
  ```

#### `dt date to-epoch`

- **Synopsis:** `dt date to-epoch [--layout <fmt>] [--ms] [--utc] <time...|stdin>`
- **Purpose:** Convert human-readable timestamps into Unix seconds (default) or milliseconds (`--ms`). Auto-detects common layouts.
- **Example:**

  ```sh
  dt date to-epoch '2025-09-17T12:34:56Z'
  # Output
  # 1758112496

  dt date to-epoch --layout '2006-01-02 15:04:05' --ms '2025-09-17 05:00:00'
  # Output
  # 1758085200000
  ```

#### `dt date from-epoch`

- **Synopsis:** `dt date from-epoch [--format rfc3339|unix|unixms|layout] [--layout <fmt>] [--utc] <epoch...|stdin>`
- **Purpose:** Convert Unix timestamps (seconds or milliseconds) into formatted time strings.
- **Example:**

  ```sh
  dt date from-epoch 1758112496 --utc
  # Output
  # 2025-09-17T12:34:56Z

  echo 1758085200000 | dt date from-epoch --format layout --layout '2006-01-02 15:04:05' --utc
  # Output
  # 2025-09-17 05:00:00
  ```

#### `dt date add`

- **Synopsis:** `dt date add --duration <GoDuration> [--from <time|epoch>] [--format rfc3339|unix|unixms|layout] [--layout <fmt>] [--utc]`
- **Purpose:** Add or subtract Go-style durations (`1h30m`, `-15m`, etc.) from either _now_ or a supplied timestamp.
- **Example:**

  ```sh
  dt date add --duration '1h30m' --from '2025-09-17T12:00:00Z'
  # Output
  # 2025-09-17T13:30:00Z

  dt date add --duration '-48h' --from 1758112496 --format unix
  # Output
  # 1757939696
  ```

### UUID Command

#### `dt uuid new`

- **Synopsis:** `dt uuid new [-n <count>]`
- **Purpose:** Generate cryptographically secure UUIDv4 identifiers.
- **Flags:**
  - `-n`, `--count` — number of UUIDs to emit (default: `1`).
- **Example:**
  ```sh
  dt uuid new -n 2
  # Output
  # 8d6b2b48-5ad7-4808-8ed1-a01a2b4dbf5b
  # c3d2c1ac-4c05-4441-83dd-99e6213d6f5a
  ```

### Environment Commands

#### `dt env from-json`

- **Synopsis:** `dt env from-json [--uppercase] [--prefix <PFX>] [--flatten] [--sep _]`
- **Purpose:** Convert a JSON object into `KEY=VALUE` lines suitable for exporting into shells or `.env` files.
- **Flag nuances:**
  - `--uppercase` — transform keys to uppercase before any prefix is applied.
  - `--prefix` — prepend a string to every key (e.g., `APP_`).
  - `--flatten` — recurse through nested objects, joining segments with `--sep` (default `_`). Arrays are serialized as compact JSON to preserve order.
  - `--sep` — separator used when flattening nested keys.
- **Example input (`config.json`):**
  ```json
  {
    "host": "db.local",
    "port": 5432,
    "auth": {
      "user": "svc",
      "scopes": ["read", "write"]
    }
  }
  ```
- **Example command:**
  ```sh
  cat config.json | dt env from-json --flatten --sep '_' --uppercase --prefix APP_
  # Output (keys sorted alphabetically)
  # APP_AUTH_SCOPES=["read","write"]
  # APP_AUTH_USER=svc
  # APP_HOST=db.local
  # APP_PORT=5432
  ```

#### `dt env from-kv`

- **Synopsis:** `dt env from-kv [--uppercase] [--prefix <PFX>]`
- **Purpose:** Convert `key: value` style pairs (e.g., copied from YAML) into `KEY=VALUE` lines. Blanks and `#` comments are ignored.
- **Example:**
  ```sh
  printf '# Service\nhost: localhost\nport: 8080\nrelease: canary\n' | dt env from-kv --uppercase
  # Output
  # HOST=localhost
  # PORT=8080
  # RELEASE=canary
  ```

### Completion Command

#### `dt completion`

- **Synopsis:** `dt completion {bash|zsh|fish|powershell}`
- **Purpose:** Emit shell completion scripts for your chosen shell. Pipe or redirect the result to the appropriate configuration file.
- **Examples:**

  ```sh
  # Bash (current shell)
  source <(dt completion bash)

  # Fish
  mkdir -p ~/.config/fish/completions
  dt completion fish > ~/.config/fish/completions/dt.fish
  ```

---

## Workflow Tips

- Chain commands freely: `dt json pretty | bat -l json`, `kubectl get pod -o json | dt json pretty`, or `terraform output -json | dt env from-json` for `.env` exports.
- When piping secret material, combine with `base64 encode --no-pad` or `env from-json --prefix` to match your deployment tooling.
- Use `dt date --help` to browse the built-in layouts before reaching for custom formatting.

---

## Development

Run tests before shipping changes:

```sh
go test ./...
```

---
