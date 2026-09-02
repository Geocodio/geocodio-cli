# geocodio-cli

A command-line interface for the [Geocodio](https://www.geocod.io/) API that lets you geocode addresses, reverse geocode coordinates, calculate distances, and process spreadsheets—all from your terminal.

Whether you're geocoding a single address or processing thousands in batch, this CLI gives you quick access to Geocodio's powerful geocoding capabilities without writing any code.

## Table of Contents

- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Cost and Billing](#cost-and-billing)
- [Commands](#commands)
  - [Geocoding](#geocoding)
  - [Reverse Geocoding](#reverse-geocoding)
  - [Distance Calculation](#distance-calculation)
  - [Distance Matrix](#distance-matrix)
  - [Async Distance Jobs](#async-distance-jobs)
  - [Spreadsheet Processing](#spreadsheet-processing)
- [Output Formats](#output-formats)
- [Global Flags](#global-flags)
- [AI Coding Assistant Skill](#ai-coding-assistant-skill)
- [Development](#development)
- [Troubleshooting](#troubleshooting)
- [License](#license)

## Getting Started

### Prerequisites

- A Geocodio API key (get one at [geocod.io](https://dash.geocod.io/apikey))

### Installation

**Quick install (recommended):**

```bash
curl -fsSL https://raw.githubusercontent.com/geocodio/geocodio-cli/main/install.sh | sh
```

**Install a specific version:**

```bash
curl -fsSL https://raw.githubusercontent.com/geocodio/geocodio-cli/main/install.sh | sh -s -- --version v1.0.0
```

**With Go:**

```bash
go install github.com/geocodio/geocodio-cli/cmd/geocodio@latest
```

### Your First Geocode

Once installed, set your API key and try geocoding an address:

```bash
export GEOCODIO_API_KEY=your-api-key
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC"
```

You'll see the latitude, longitude, and formatted address returned by the API.

## Configuration

### API Key

You can provide your API key in two ways:

**Environment variable (recommended):**

```bash
export GEOCODIO_API_KEY=your-api-key
```

**Command-line flag:**

```bash
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --api-key your-api-key
```

> [!TIP]
> Add the export command to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) so you don't have to set it every session.

### API Key Permissions

Geocodio API keys carry per-feature permissions. The `lists` and `distance` commands (`distance`, `distance-matrix`, and `distance-jobs`) only work when spreadsheet and distance access are enabled for the key you are using. You can turn them on at [dash.geocod.io/apikey](https://dash.geocod.io/apikey).

If those commands return a `403` while `geocode` works fine, the key is valid and the permission is missing.

## Cost and Billing

Every command here spends lookups on your Geocodio account. One lookup is one credit, and a run costs:

```
records × (1 + number of append categories)
```

- **Appends are not free.** Each category in `--fields` is a full extra lookup per record. `--fields timezone,cd` on 1,000 addresses costs 3,000 lookups, not 1,000.
- **Distance commands** cost geocoding for any input given as an address, plus `origins × destinations × mode multiplier`, plus any appends. Straightline counts 1× per pair, driving counts 2×.
- **Free tier:** the first 2,500 lookups each day are free on pay-as-you-go and Flex plans. The allowance resets daily and does not roll over.
- **Zero-result lookups are not billed.** An address that returns no match costs nothing.

### Ways to Spend Less

- Pass `--skip-geocoding` when you already have coordinates and only need field appends.
- Store the `stable_address_key` returned by `--show-address-key`, then skip records you have geocoded before.
- Pass coordinates instead of addresses to `distance`, `distance-matrix`, and `distance-jobs`, so there is nothing left to geocode.
- Use `--mode straightline` unless you need real driving routes. Driving costs twice as much per pair.
- Request only the append fields you use. Every extra category multiplies the whole run.

### Set a Limit Before a Large Run

A usage limit on your account caps what a run can spend, which is worth setting before the first big file: [Set a usage limit](https://www.geocod.io/guides/set-a-usage-limit).

For plan and per-lookup rates, see [Geocodio pricing](https://www.geocod.io/pricing). To size a plan against your own volume, use the [plan calculator](https://www.geocod.io/find-my-plan).

## Commands

### Geocoding

Convert addresses into geographic coordinates.

**Single address:**

```bash
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC"
```

**With data appends:**

Geocodio can return additional data like timezone, congressional district, census data, and more. Specify fields with the `--fields` flag:

```bash
geocodio geocode "30 Rockefeller Plaza, New York NY" --fields timezone,cd
```

> [!IMPORTANT]
> Appends are not free. Each `--fields` category is a full extra lookup per record, so `--fields timezone,cd` on 1,000 addresses costs 3,000 lookups instead of 1,000. See [Cost and Billing](#cost-and-billing) and [Geocodio's field documentation](https://www.geocod.io/docs/#data-appends-fields).

United Kingdom data appends are available too — for example Westminster and devolved constituencies and local authority wards:

```bash
geocodio geocode "10 Downing St, London" --country "United Kingdom" --fields uk-westminster,uk-local
```

UK field appends: `uk-westminster`, `uk-westminster-next`, `uk-devolved`, `uk-devolved-next`, `uk-local`, `uk-local-next`. The `-next` variants return upcoming boundary changes.

> [!NOTE]
> United Kingdom addresses need a Flex or Unlimited+UK plan. UK geocoding is not available on pay-as-you-go. See [Geocodio pricing](https://www.geocod.io/pricing).

**Batch geocoding from a file:**

For processing many addresses at once, create a file with one address per line and use the `--batch` flag:

```bash
geocodio geocode --batch addresses.txt
```

When running in a terminal, you'll see a spinner while the batch processes.

> [!IMPORTANT]
> `--batch` caps at **10,000 lookups per request**, and appends count toward that cap: 5,000 addresses with one `--fields` category already reaches the limit. A full 10,000-lookup batch can take around 600 seconds, so set script and client timeouts accordingly.
>
> For anything larger, and for any spreadsheet, prefer [`lists upload`](#spreadsheet-processing). It runs asynchronously, handles up to 10 million lookups, needs no file splitting, and reports progress with `--watch`.

**With inline distance calculations:**

Geocode an address and calculate distances to one or more destinations in a single request:

```bash
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --destinations "New York" --destinations "Boston"
```

With driving mode for actual driving distance and duration:

```bash
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --destinations "New York" --distance-mode driving
```

**With stable address key:**

Show the stable address key for a result, which can be used in future requests instead of an address:

```bash
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --show-address-key
```

**JSON output:**

```bash
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --json
```

**All geocode flags:**

| Flag | Alias | Description |
|------|-------|-------------|
| `--batch` | `-b` | File containing addresses (one per line) |
| `--fields` | `-f` | Data append fields (comma-separated) |
| `--limit` | `-l` | Maximum number of results per address |
| `--country` | `-c` | Country hint, recommended for non-US data (e.g. `USA`, `Canada`, `United Kingdom`) |
| `--destinations` | `-d` | Destination addresses or coordinates for distance calculation (repeatable) |
| `--distance-mode` | `-m` | Distance mode: `driving` or `straightline` |
| `--distance-units` | `-u` | Distance units: `miles` or `km` |
| `--distance-max-distance` | `--distance-radius` | Radius limit: only keep destinations within this distance |
| `--distance-min-distance` | | Only keep destinations at least this far away |
| `--distance-max-duration` | | Only keep destinations within this travel time in seconds (driving only) |
| `--distance-min-duration` | | Only keep destinations at least this many seconds away (driving only) |
| `--distance-max-results` | | Only keep the N nearest destinations per result |
| `--distance-order-by` | | Sort destinations by: `distance` or `duration` |
| `--distance-sort-order` | | Sort direction: `asc` or `desc` |
| `--show-address-key` | | Show stable address key in output |

> [!TIP]
> Pass `--country` whenever your data is not from the United States. Without it the API falls back to the US, which returns confident but wrong matches for addresses elsewhere.

### Reverse Geocoding

Convert coordinates back into addresses.

**Single coordinate pair:**

```bash
geocodio reverse "38.8976,-77.0365"
```

**Batch reverse geocoding:**

```bash
geocodio reverse --batch coordinates.txt
```

The coordinates file should have one `lat,lng` pair per line:

```
38.8976,-77.0365
40.7128,-74.0060
34.0522,-118.2437
```

**Skip geocoding (fields only):**

Get only field appends (timezone, census data, etc.) for coordinates without reverse geocoding an address:

```bash
geocodio reverse "40.7588,-73.9788" --skip-geocoding --fields timezone,cd
```

**With inline distance calculations:**

```bash
geocodio reverse "38.8976,-77.0365" --destinations "New York" --distance-mode driving
```

**All reverse flags:**

| Flag | Alias | Description |
|------|-------|-------------|
| `--batch` | `-b` | File containing coordinates (one per line) |
| `--fields` | `-f` | Data append fields (comma-separated) |
| `--limit` | `-l` | Maximum number of results per coordinate |
| `--skip-geocoding` | | Skip reverse geocoding, only return field appends |
| `--destinations` | `-d` | Destination addresses or coordinates for distance calculation (repeatable) |
| `--distance-mode` | `-m` | Distance mode: `driving` or `straightline` |
| `--distance-units` | `-u` | Distance units: `miles` or `km` |
| `--distance-max-distance` | `--distance-radius` | Radius limit: only keep destinations within this distance |
| `--distance-min-distance` | | Only keep destinations at least this far away |
| `--distance-max-duration` | | Only keep destinations within this travel time in seconds (driving only) |
| `--distance-min-duration` | | Only keep destinations at least this many seconds away (driving only) |
| `--distance-max-results` | | Only keep the N nearest destinations per result |
| `--distance-order-by` | | Sort destinations by: `distance` or `duration` |
| `--distance-sort-order` | | Sort direction: `asc` or `desc` |
| `--show-address-key` | | Show stable address key in output |

### Distance Calculation

Calculate distances from one origin to one or more destinations.

**Basic usage:**

```bash
geocodio distance "Washington DC" "New York"
```

**Multiple destinations:**

```bash
geocodio distance "Washington DC" "New York" "Boston" "Philadelphia"
```

**With options:**

```bash
geocodio distance "Washington DC" "New York" --mode driving --units km
```

**Radius limiting:**

Use `--radius` (an alias for `--max-distance`) to drop destinations that fall outside a radius around the origin. The value is interpreted in whatever `--units` is set to.

```bash
# Only the destinations within 150 miles of Washington DC
geocodio distance "Washington DC" "New York" "Boston" "Philadelphia" --radius 150

# Same, in kilometers
geocodio distance "Washington DC" "New York" "Boston" "Philadelphia" --radius 250 --units km

# A ring: at least 50 miles out, at most 150
geocodio distance "Washington DC" "New York" "Boston" "Philadelphia" --min-distance 50 --max-distance 150
```

You can also limit by travel time instead of distance, and keep only the nearest few:

```bash
# Within a two-hour drive
geocodio distance "Washington DC" "New York" "Boston" "Philadelphia" --mode driving --max-duration 7200

# The three closest destinations, longest drive first
geocodio distance "Washington DC" "New York" "Boston" "Philadelphia" \
  --max-results 3 --order-by duration --sort-order desc
```

**Distance flags:**

| Flag | Alias | Default | Description |
|------|-------|---------|-------------|
| `--mode` | `-m` | `driving` | Routing mode: `driving` or `straightline` |
| `--units` | `-u` | `miles` | Distance units: `miles` or `km` |
| `--country` | `-c` | | Country to append to addresses, recommended for non-US data (e.g. `USA`, `Canada`, `United Kingdom`) |
| `--max-distance` | `--radius` | | Radius limit: only keep destinations within this distance (in `--units`) |
| `--min-distance` | | | Only keep destinations at least this far away (in `--units`) |
| `--max-duration` | | | Only keep destinations within this travel time in seconds (driving mode only) |
| `--min-duration` | | | Only keep destinations at least this many seconds away (driving mode only) |
| `--max-results` | | | Only keep the N nearest destinations |
| `--order-by` | | `distance` | Sort destinations by: `distance` or `duration` |
| `--sort-order` | | `asc` | Sort direction: `asc` or `desc` |

> [!TIP]
> Use `straightline` mode for quick "as the crow flies" distances when you don't need actual driving routes.

> [!NOTE]
> `--max-duration` and `--min-duration` only apply in `driving` mode, since `straightline` results have no travel time.

> [!IMPORTANT]
> A distance run costs geocoding for every input given as an address, plus one lookup per origin/destination pair in `straightline` mode or two per pair in `driving` mode, plus any appends. Passing coordinates instead of addresses removes the geocoding cost. See [Cost and Billing](#cost-and-billing).

### Distance Matrix

Calculate distances between multiple origins and multiple destinations. This is useful for logistics, delivery routing, or finding the closest locations.

```bash
geocodio distance-matrix --origins origins.txt --destinations destinations.txt
```

Both files should contain one location per line (addresses or coordinates).

The same radius and travel-time limits available on `distance` apply here, per origin. This is how you answer "which of my stores is within 10 miles of each customer?" without post-processing the full matrix:

```bash
geocodio distance-matrix --origins customers.txt --destinations stores.txt --radius 10

# Nearest store to each customer, by driving time
geocodio distance-matrix --origins customers.txt --destinations stores.txt \
  --mode driving --max-results 1 --order-by duration
```

**Distance matrix flags:**

| Flag | Alias | Required | Default | Description |
|------|-------|----------|---------|-------------|
| `--origins` | `-o` | Yes | — | File containing origin locations |
| `--destinations` | `-d` | Yes | — | File containing destination locations |
| `--mode` | `-m` | No | `driving` | Routing mode: `driving` or `straightline` |
| `--units` | `-u` | No | `miles` | Distance units: `miles` or `km` |
| `--country` | `-c` | No | | Country to append to addresses, recommended for non-US data |
| `--max-distance` | `--radius` | No | | Radius limit: only keep destinations within this distance (in `--units`) |
| `--min-distance` | | No | | Only keep destinations at least this far away (in `--units`) |
| `--max-duration` | | No | | Only keep destinations within this travel time in seconds (driving mode only) |
| `--min-duration` | | No | | Only keep destinations at least this many seconds away (driving mode only) |
| `--max-results` | | No | | Only keep the N nearest destinations per origin |
| `--order-by` | | No | `distance` | Sort destinations by: `distance` or `duration` |
| `--sort-order` | | No | `asc` | Sort direction: `asc` or `desc` |

> [!IMPORTANT]
> A matrix costs `origins × destinations` pairs, doubled in `driving` mode, plus geocoding for any inputs given as addresses. A 500 × 500 matrix is 250,000 pairs, so check the numbers (and your [usage limit](https://www.geocod.io/guides/set-a-usage-limit)) before a large run. See [Cost and Billing](#cost-and-billing).

### Async Distance Jobs

For large distance calculations, use async jobs. These run in the background and you can check their status or download results later.

**Create a job:**

```bash
geocodio distance-jobs create --name "My Job" --origins origins.txt --destinations destinations.txt
```

**Create and watch progress:**

```bash
geocodio distance-jobs create --name "My Job" --origins origins.txt --destinations destinations.txt --watch
```

When running in a terminal, you'll see an animated progress bar showing the job's completion status.

**List all jobs:**

```bash
geocodio distance-jobs list
```

**Check job status:**

```bash
geocodio distance-jobs status 12345

# Watch until completion
geocodio distance-jobs status 12345 --watch
```

**Download results:**

```bash
# Output to stdout
geocodio distance-jobs download 12345

# Save to file
geocodio distance-jobs download 12345 --output results.csv
```

**Delete a job:**

```bash
geocodio distance-jobs delete 12345
```

**Distance jobs create flags:**

| Flag | Alias | Required | Default | Description |
|------|-------|----------|---------|-------------|
| `--name` | `-n` | Yes | — | Job name for identification |
| `--origins` | `-o` | Yes | — | File containing origin locations |
| `--destinations` | `-d` | Yes | — | File containing destination locations |
| `--mode` | `-m` | No | `driving` | Routing mode: `driving` or `straightline` |
| `--units` | `-u` | No | `miles` | Distance units: `miles` or `km` |
| `--watch` | `-w` | No | `false` | Watch progress until completion |

> [!IMPORTANT]
> Async jobs are priced the same way as `distance-matrix`: `origins × destinations` pairs, doubled in `driving` mode, plus geocoding for address inputs. The job size is the whole matrix, so a job created by mistake can spend a lot of lookups. See [Cost and Billing](#cost-and-billing).

### Spreadsheet Processing

Geocode CSV or Excel files without splitting them. Geocodio processes the file asynchronously and returns results with coordinates appended to your data. This is the recommended path for spreadsheets and for anything above the 10,000-lookup `--batch` cap: a list handles up to 10 million lookups.

**Upload a file:**

```bash
geocodio lists upload data.csv --direction forward --format "{{A}} {{B}}, {{C}}"
```

The `--format` flag tells Geocodio which columns contain address components. Use `{{A}}`, `{{B}}`, `{{C}}`, etc. to reference columns:

- `{{A}}` = Column A (first column)
- `{{B}}` = Column B (second column)
- And so on...

**Example:** If your CSV has street in column A, city in column B, and state in column C:

```bash
geocodio lists upload addresses.csv --direction forward --format "{{A}}, {{B}}, {{C}}"
```

**Upload and watch progress:**

```bash
geocodio lists upload data.csv --direction forward --format "{{A}}" --watch
```

**List all uploaded spreadsheets:**

```bash
geocodio lists list
```

**Check status:**

```bash
geocodio lists status 12345

# Watch until completion
geocodio lists status 12345 --watch
```

**Download results:**

```bash
# Output to stdout
geocodio lists download 12345

# Save to file
geocodio lists download 12345 --output geocoded.csv
```

**Delete a spreadsheet:**

```bash
geocodio lists delete 12345
```

**Lists upload flags:**

| Flag | Alias | Required | Description |
|------|-------|----------|-------------|
| `--direction` | `-d` | Yes | `forward` (address→coords) or `reverse` (coords→address) |
| `--format` | `-f` | Yes | Column format template (e.g., `{{A}} {{B}}, {{C}}`) |
| `--watch` | `-w` | No | Watch progress until completion |
| `--callback` | — | No | URL to receive a POST request when processing completes |
| `--fields` | `-F` | No | Data append fields (comma-separated, e.g., `cd,timezone`) |

> [!WARNING]
> Large spreadsheets can take time to process. Use the `--watch` flag or check status periodically rather than waiting for immediate results.

> [!IMPORTANT]
> A list costs `records × (1 + number of append categories)` lookups, the same as batch geocoding. `--fields cd,timezone` on a 100,000-row spreadsheet is 300,000 lookups. Spreadsheet access also has to be enabled on your API key. See [Cost and Billing](#cost-and-billing).

## Output Formats

The CLI supports multiple output formats to fit different workflows.

### Human-Readable (Default)

The default output is formatted for easy reading in your terminal. When connected to a terminal, you'll see:

- **Colored output** with status indicators (green for completed, yellow for processing, red for failed)
- **Styled labels** for better visual hierarchy
- **Progress indicators** during batch operations and watch mode

### JSON Output

For scripting and programmatic use, get raw JSON with the `--json` flag:

```bash
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --json
```

This returns the complete API response, perfect for piping to `jq` or processing in scripts.

### Agent Output (Markdown)

For AI assistants and LLMs, use the `--agent` flag to get clean markdown tables:

```bash
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --agent
```

This outputs structured markdown that's easy for language models to parse:

```markdown
## Geocode Result

| Field | Value |
|-------|-------|
| Matched Address | 1600 Pennsylvania Ave NW, Washington, DC 20500 |
| Coordinates | 38.8976763, -77.0365298 |
| Accuracy | rooftop (1.00) |
```

### Disabling Colors

Colors are automatically disabled when output is piped or redirected. To explicitly disable colors:

```bash
# Using the flag
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --no-color

# Using the environment variable
NO_COLOR=1 geocodio geocode "1600 Pennsylvania Ave NW, Washington DC"
```

> [!TIP]
> The `NO_COLOR` environment variable follows the [no-color.org](https://no-color.org) standard and works across many CLI tools.

## Global Flags

These flags work with all commands:

| Flag | Description |
|------|-------------|
| `--api-key` | Geocodio API key (or use `GEOCODIO_API_KEY` env var) |
| `--json` | Output results as JSON |
| `--agent` | Output results as markdown (for LLM consumption) |
| `--no-color` | Disable colored output (also respects `NO_COLOR` env var) |
| `--base-url` | Override API base URL (useful for testing or enterprise endpoints) |
| `--debug` | Show HTTP request/response details |
| `--version` | Show version information |
| `--help` | Show help for any command |

## AI Coding Assistant Skill

The Geocodio CLI includes a skill that teaches AI coding assistants how to use it. Once installed, your assistant can geocode addresses, reverse geocode coordinates, calculate distances, and process spreadsheets on your behalf.

### Compatibility

Works with any AI coding assistant that supports skills, including [Claude Code](https://claude.ai/code), [Cursor](https://cursor.com), [Amp](https://amp.dev), and [Codex](https://openai.com/codex).

### Installation

```bash
npx skills add geocodio/geocodio-cli
```

### What's Included

The skill teaches your assistant:

- All CLI commands (geocode, reverse, distance, distance-matrix, distance-jobs, lists)
- Output format flags (`--json`, `--agent`) and when to use each
- Batch processing workflows and file size limits
- Data append fields (timezone, congressional district, census, etc.)
- Common multi-step workflows like geocoding a CSV and downloading results

### Usage

Once installed, the skill activates automatically when you ask your assistant to geocode, look up an address, calculate distances, or process a spreadsheet. No special syntax needed.

**Example prompts:**

- "Geocode this list of addresses and save the results as JSON"
- "What congressional district is 1600 Pennsylvania Ave in?"
- "Calculate driving distances from our warehouse to these 50 customers"
- "Upload this CSV and geocode the addresses in columns B and C"

### Cost Guidance for Agents

Before you let an assistant run large jobs on your account, point it at [Cost and Billing](#cost-and-billing). It covers how lookups are counted, why appends multiply a run, and when `lists upload` beats `--batch`, which is the context an agent needs to avoid an expensive surprise.

### Skill Source

The skill definition is in [`skills/geocodio/SKILL.md`](skills/geocodio/SKILL.md).

## Development

### Building

```bash
# Build binary to ./bin/geocodio
make build

# Install to your $GOPATH/bin
make install
```

### Testing

Tests use [go-vcr](https://github.com/dnaeon/go-vcr) to record and replay HTTP interactions. This lets tests run without hitting the live API.

**Run tests:**

```bash
make test
```

**Run tests with coverage:**

```bash
make cover
```

This generates a coverage report at `coverage.html`.

#### VCR Cassettes

Tests replay recorded API responses from "cassette" files stored in `internal/api/testdata/`. This approach:

- Makes tests fast and deterministic
- Avoids API rate limits and costs during development
- Allows tests to run without an API key

**Recording new cassettes:**

When adding new tests or updating existing ones, you'll need to record fresh API interactions:

```bash
GEOCODIO_API_KEY=your-real-key make record-cassettes
```

> [!IMPORTANT]
> API keys are automatically redacted from recorded cassettes, so they're safe to commit.

**How VCR works in this project:**

1. When `VCR_MODE=record` is set, tests make real API calls and save responses to YAML cassette files
2. During normal test runs, requests are matched against recorded cassettes and replayed
3. The matcher ignores API keys, so cassettes work regardless of which key was used during recording
4. If a cassette doesn't exist and `VCR_MODE` isn't set to `record`, the test is skipped

**Adding a new test with VCR:**

```go
func TestNewEndpoint(t *testing.T) {
    client := newTestClient(t, "new_endpoint")

    // Make your API call
    resp, err := client.SomeMethod(ctx, req)

    // Assert results
    require.NoError(t, err)
    assert.Equal(t, expected, resp)
}
```

Then record the cassette:

```bash
GEOCODIO_API_KEY=your-key VCR_MODE=record go test -v -run TestNewEndpoint ./internal/api/...
```

### Smoke Testing

The unit tests use recorded VCR cassettes and don't hit the live API. To verify that the CLI works end-to-end against the real Geocodio API, run the smoke test:

```bash
GEOCODIO_API_KEY=your-real-key ./scripts/smoke-test.sh
```

This exercises every command (geocode, reverse, distance, distance-matrix, distance-jobs, lists) with various flag combinations, output formats, and error cases. It also runs lifecycle tests that create and delete distance jobs and spreadsheet uploads.

> [!NOTE]
> The smoke test makes real API calls, which may count against your usage. It creates temporary resources (distance jobs, spreadsheet uploads) that are automatically deleted during the test run.

### Linting

```bash
make lint
```

Requires [golangci-lint](https://golangci-lint.run/usage/install/).

### Cleaning Up

```bash
make clean
```

## Troubleshooting

### "API key required"

You haven't set your API key. Either:

```bash
export GEOCODIO_API_KEY=your-api-key
```

Or pass it directly:

```bash
geocodio geocode "address" --api-key your-api-key
```

### "batch size exceeds maximum of 10,000"

Geocodio's batch endpoints accept a maximum of 10,000 lookups per request, and appends count toward that cap: 5,000 addresses with one `--fields` category reaches it.

Above that limit, upload the file as a list instead of splitting it. A list runs asynchronously, takes up to 10 million lookups, and reports progress:

```bash
geocodio lists upload large_file.csv --direction forward --format "{{A}}" --watch
geocodio lists download 12345 --output geocoded.csv
```

Splitting the file still works if you need results in a single stream, but it costs the same lookups and takes longer end to end:

```bash
split -l 10000 large_file.txt chunk_
for f in chunk_*; do geocodio geocode --batch "$f"; done
```

### "geocodio API error (403)" from `lists` or `distance` commands

A `403` from these commands usually means the API key is fine but the permission is not enabled. Spreadsheet and distance access are per-key settings at [dash.geocod.io/apikey](https://dash.geocod.io/apikey). If `geocodio geocode "1600 Pennsylvania Ave NW, Washington DC"` succeeds with the same key, the key itself is valid.

### "invalid coordinate format"

Reverse geocoding expects coordinates in `lat,lng` format:

```bash
# Correct
geocodio reverse "38.8976,-77.0365"

# Incorrect
geocodio reverse "38.8976 -77.0365"  # Missing comma
geocodio reverse "-77.0365,38.8976"  # Longitude first (should be latitude)
```

### Debugging API Issues

Use the `--debug` flag to see the full HTTP request and response:

```bash
geocodio geocode "1600 Pennsylvania Ave" --debug
```

### Colors Not Displaying

If you're not seeing colored output:

1. Make sure you're running in a terminal (not piping output)
2. Check that `NO_COLOR` isn't set in your environment
3. Try using `FORCE_COLOR=1` to force color output

## Migrating from v1.x

If you're upgrading from the previous Geocodio CLI (v1.x), here are the breaking changes:

| v1.x | v2.x | Notes |
|------|------|-------|
| `geocodio create file.csv "{{A}}"` | `geocodio lists upload file.csv --direction forward --format "{{A}}"` | `--direction` is now required (`forward` or `reverse`) |
| `geocodio status 123` | `geocodio lists status 123` | Commands are grouped under `lists` |
| `geocodio download 123 > out.csv` | `geocodio lists download 123 --output out.csv` | Use `--output` flag instead of shell redirect (stdout still works too) |
| `geocodio remove 123` | `geocodio lists delete 123` | Renamed for consistency |
| `geocodio list` | `geocodio lists list` | Grouped under `lists` |
| `--follow` | `--watch` | Flag renamed |
| `--apikey` / `-k` | `--api-key` | Flag renamed (env var `GEOCODIO_API_KEY` unchanged) |
| `--hostname` / `-n` | `--base-url` | Flag renamed |

## License

MIT
