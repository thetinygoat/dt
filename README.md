# dt - Dev toolkit that doesn't suck 🛠️

Tired of juggling 15 different CLI tools just to format some JSON and convert a timestamp? Yeah, me too. 

`dt` is a single Go binary that packs all those everyday dev utilities you actually use into one fast, pipe-friendly command. JSON formatting, date wrangling, UUID generation, Base64 encoding, hash generation, and environment variable helpers - all in one place.

---

## Getting Started

```sh
go install ./...
# or build it yourself
go build -o dt .
```

You'll need Go 1.21+. First run automatically grabs the Cobra dependency.

---

## How it works

`dt` is designed to play nice with Unix pipes. It'll read from stdin when you pipe data to it, otherwise it uses whatever arguments you pass. Results go to stdout, errors to stderr, so you can chain commands without worrying about breaking your pipeline.

Hit `dt --help` for the full command list, or `dt <command> --help` for specific flags.

### Quick examples

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

## What can it do? 

Here's everything `dt` can handle. Examples use Unix-style shells - just swap in your own values.

### JSON Commands

#### `dt json pretty`

Makes ugly JSON readable. It's smart enough to unwrap stringified JSON up to 3 levels deep, so you can pipe those gnarly escaped payloads straight in.

- **Usage:** `dt json pretty [--indent 2]`
- **Flags:**
  - `--indent <n>` - spaces per level (default: 2)
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

Need to embed JSON in config files or code? This escapes it properly so you don't have to manually escape all those quotes.

- **Usage:** `dt json stringify [--compact] [--no-quotes] <json|stdin>`
- **Flags:**
  - `--compact` - removes all whitespace for the smallest output
  - `--no-quotes` - skip the outer quotes (handy when your tooling adds them)
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

Standard Base64 encoding with options for URL-safe output and padding control.

- **Usage:** `dt base64 encode [--url] [--no-pad]`
- **Flags:**
  - `--url` - use URL-safe characters (`-` and `_` instead of `+` and `/`)
  - `--no-pad` - drop the `=` padding (great for JWTs)
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

Decodes Base64 back to the original. Smart enough to handle both padded and unpadded inputs.

- **Usage:** `dt base64 decode [--url]`
- **Flags:**
  - `--url` - decode URL-safe Base64
- **Example:**
  ```sh
  echo 'YXBwOnNlY3JldA==' | dt base64 decode
  # Output
  # app:secret
  ```

### Hash Commands

Generate hashes with all the common algorithms. Supports salting and multiple output formats.

#### `dt hash <algorithm>`

- **Usage:** `dt hash sha256|sha512|sha3-256|sha3-512|sha1|md5 [--encoding hex|base64] [--salt <value>]`
- **Flags:**
  - `--encoding` - hex (default) or base64 output
  - `--salt` - add a salt string before hashing
- **Example:**

  ```sh
  echo -n 'hello' | dt hash sha3-256
  # 3338be694f50c5f338814986cdf0686453a888b84f424d792af4b9202398f392

  dt hash sha256 --encoding base64 hello
  # LPJNul+wow4m6DsqxbninhsWHlwfp0JecwQzYpOLmCQ=

  echo -n 'hello' | dt hash md5 --salt pepper
  # 6967321c83e9f01a33e7edecce748877
  ```

### Text Commands

#### `dt text join`

Perfect for turning spreadsheet columns or lists into SQL IN clauses, shell arrays, or any delimited format you need.

- **Usage:** `dt text join [--sep ","] [--quote single|double|none] [--split lines|tab|csv] [--trim] [--skip-empty] [--unique] [items...]`
- **Flags:**
  - `--sep` - what to put between items (supports `\n`, `\t`, `\r`, `\0`). Default: `,`
  - `--quote` - single quotes (default), double quotes, or none
  - `--split` - how to parse input: `lines`, `tab`, or `csv` (handles quoted commas)
  - `--trim` / `--no-trim` - strip whitespace from each item (on by default)
  - `--skip-empty` - ignore empty items after trimming (on by default)  
  - `--unique` - remove duplicates, keeping the first occurrence
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

All the date/time wrangling you need. Commands work with arguments or piped input. Layout strings use Go's reference time format (`2006-01-02 15:04:05`).

#### `dt date now`

Get the current time in whatever format you need.

- **Usage:** `dt date now [--format rfc3339|unix|unixms|layout] [--layout <fmt>] [--utc]`
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

Turn readable timestamps into Unix epochs. Auto-detects most common formats.

- **Usage:** `dt date to-epoch [--layout <fmt>] [--ms] [--utc] <time...|stdin>`
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

Convert those Unix timestamps back into something humans can read.

- **Usage:** `dt date from-epoch [--format rfc3339|unix|unixms|layout] [--layout <fmt>] [--utc] <epoch...|stdin>`
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

Time math made easy. Add or subtract durations using Go's format (`1h30m`, `-15m`, etc.).

- **Usage:** `dt date add --duration <GoDuration> [--from <time|epoch>] [--format rfc3339|unix|unixms|layout] [--layout <fmt>] [--utc]`
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

Generate UUIDv4s. Cryptographically secure, as many as you need.

- **Usage:** `dt uuid new [-n <count>]`
- **Flags:**
  - `-n`, `--count` - how many UUIDs to generate (default: 1)
- **Example:**
  ```sh
  dt uuid new -n 2
  # Output
  # 8d6b2b48-5ad7-4808-8ed1-a01a2b4dbf5b
  # c3d2c1ac-4c05-4441-83dd-99e6213d6f5a
  ```

### Environment Commands

#### `dt env from-json`

Turn JSON configs into shell environment variables. Great for Docker, CI/CD, or any `.env` workflow.

- **Usage:** `dt env from-json [--uppercase] [--prefix <PFX>] [--flatten] [--sep _]`
- **Flags:**
  - `--uppercase` - MAKE_KEYS_LIKE_THIS
  - `--prefix` - add a prefix to every key (like `APP_`)
  - `--flatten` - turn nested objects into flat keys with separators
  - `--sep` - what to use for separating nested keys (default: `_`)
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

Parse `key: value` pairs (like from YAML) into environment variables. Ignores comments and blank lines.

- **Usage:** `dt env from-kv [--uppercase] [--prefix <PFX>]`
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

Get tab completion working in your shell. Because nobody likes typing long command names.

- **Usage:** `dt completion {bash|zsh|fish|powershell}`
- **Examples:**

  ```sh
  # Bash (current shell)
  source <(dt completion bash)

  # Fish
  mkdir -p ~/.config/fish/completions
  dt completion fish > ~/.config/fish/completions/dt.fish
  ```

---

## Pro tips 💡

`dt` plays well with other tools. Try chaining commands like:
- `dt json pretty | bat -l json` for syntax highlighted output  
- `kubectl get pod -o json | dt json pretty` to make K8s output readable
- `terraform output -json | dt env from-json` to generate `.env` files

When handling secrets, `dt base64 encode --no-pad` or `dt env from-json --prefix` can help match your deployment tooling's format.

---

## Contributing

Want to hack on `dt`? Cool! Just make sure the tests pass:

```sh
go test ./...
```
