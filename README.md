# geocodio-cli

A command-line interface for the [Geocodio](https://www.geocod.io/) API that lets you geocode addresses, reverse geocode coordinates, calculate distances, and process spreadsheets—all from your terminal.

Whether you're geocoding a single address or processing thousands in batch, this CLI gives you quick access to Geocodio's powerful geocoding capabilities without writing any code.

## Table of Contents

- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Commands](#commands)
  - [Geocoding](#geocoding)
  - [Reverse Geocoding](#reverse-geocoding)
  - [Distance Calculation](#distance-calculation)
  - [Distance Matrix](#distance-matrix)
  - [Async Distance Jobs](#async-distance-jobs)
  - [Spreadsheet Processing](#spreadsheet-processing)
- [Output Formats](#output-formats)
- [Global Flags](#global-flags)
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
geocodio geocode "1600 Pennsylvania Ave NW" --api-key your-api-key
```

> [!TIP]
> Add the export command to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) so you don't have to set it every session.

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
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --fields timezone,cd
```

> [!NOTE]
> Data append fields may incur additional API costs. See [Geocodio's documentation](https://www.geocod.io/docs/#data-appends-fields) for available fields and pricing.

**Batch geocoding from a file:**

For processing many addresses at once, create a file with one address per line and use the `--batch` flag:

```bash
geocodio geocode --batch addresses.txt
```

When running in a terminal, you'll see a spinner while the batch processes.

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
| `--country` | `-c` | Country hint (`US` or `CA`) |
| `--destinations` | `-d` | Destination addresses or coordinates for distance calculation (repeatable) |
| `--distance-mode` | `-m` | Distance mode: `driving` or `straightline` |
| `--distance-units` | `-u` | Distance units: `miles` or `km` |
| `--show-address-key` | | Show stable address key in output |

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
geocodio reverse "38.8976,-77.0365" --skip-geocoding --fields timezone,cd
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

**Distance flags:**

| Flag | Alias | Default | Description |
|------|-------|---------|-------------|
| `--mode` | `-m` | `driving` | Routing mode: `driving` or `straightline` |
| `--units` | `-u` | `miles` | Distance units: `miles` or `km` |

> [!TIP]
> Use `straightline` mode for quick "as the crow flies" distances when you don't need actual driving routes.

### Distance Matrix

Calculate distances between multiple origins and multiple destinations. This is useful for logistics, delivery routing, or finding the closest locations.

```bash
geocodio distance-matrix --origins origins.txt --destinations destinations.txt
```

Both files should contain one location per line (addresses or coordinates).

**Distance matrix flags:**

| Flag | Alias | Required | Default | Description |
|------|-------|----------|---------|-------------|
| `--origins` | `-o` | Yes | — | File containing origin locations |
| `--destinations` | `-d` | Yes | — | File containing destination locations |
| `--mode` | `-m` | No | `driving` | Routing mode: `driving` or `straightline` |
| `--units` | `-u` | No | `miles` | Distance units: `miles` or `km` |

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

### Spreadsheet Processing

Upload CSV or Excel files for batch geocoding. Geocodio processes the file asynchronously and returns results with coordinates appended to your data.

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
geocodio geocode "1600 Pennsylvania Ave NW" --json
```

This returns the complete API response, perfect for piping to `jq` or processing in scripts.

### Agent Output (Markdown)

For AI assistants and LLMs, use the `--agent` flag to get clean markdown tables:

```bash
geocodio geocode "1600 Pennsylvania Ave NW" --agent
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
geocodio geocode "1600 Pennsylvania Ave NW" --no-color

# Using the environment variable
NO_COLOR=1 geocodio geocode "1600 Pennsylvania Ave NW"
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

## Development

See [DECISIONS.md](DECISIONS.md) for internal design decisions, including which API features are intentionally not exposed in the CLI and why.

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

Geocodio's batch endpoints accept a maximum of 10,000 items per request. Split your file into smaller chunks:

```bash
split -l 10000 large_file.txt chunk_
for f in chunk_*; do geocodio geocode --batch "$f"; done
```

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
