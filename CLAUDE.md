# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`geocodio-cli` is a command-line interface for the [Geocodio](https://www.geocod.io/) API. It geocodes addresses, reverse geocodes coordinates, calculates distances and travel times, runs N×M distance matrices and async distance jobs, and manages spreadsheet (list) geocoding jobs — all from the terminal.

Written in Go (module `github.com/geocodio/geocodio-cli`, Go 1.24+). The binary is named `geocodio`. CLI parsing uses **urfave/cli/v3**; interactive progress uses **Bubble Tea** (charmbracelet `bubbletea`/`bubbles`/`lipgloss`).

> History: this repo was rewritten from scratch. The previous `urfave/cli/v2`, command-per-package CLI now lives in the archived `geocodio-cli-legacy` repo and shares no history with this codebase.

## Build and Development Commands

The `Makefile` is the entry point for all common tasks:

```bash
make build              # Build to bin/geocodio (version injected via ldflags)
make install            # go install with version ldflags
make test               # go test -race ./...  (all packages)
make lint               # golangci-lint run
make cover              # Coverage report to coverage.html
make clean              # Remove bin/
make record-cassettes   # Re-record go-vcr cassettes (requires GEOCODIO_API_KEY)
```

Run a single test directly:

```bash
go test ./internal/cli -run TestGeocode
go test ./internal/api -run TestClient_Geocode
```

Releases are produced via GoReleaser/CI on tag push (configured outside the Makefile). `make build` derives the version from `git describe`.

## Running

```bash
export GEOCODIO_API_KEY="your_api_key"
./bin/geocodio geocode "1600 Pennsylvania Ave NW, Washington DC"

# Or pass the key explicitly
./bin/geocodio --api-key=YOUR_KEY reverse "38.9,-77.03"
```

## Architecture

The entry point is `cmd/geocodio/main.go`, which sets up signal handling (Ctrl-C → context cancel) and calls `cli.Run(ctx, os.Args)`. All logic lives under `internal/`:

| Package | Responsibility |
|---|---|
| `internal/cli` | Command definitions and action handlers. `root.go` assembles the root command and global flags. |
| `internal/api` | HTTP client (`client.go`) and typed request/response models (`types.go`), one file per endpoint group. |
| `internal/config` | Resolves API key, base URL, and debug flag from flags + environment. |
| `internal/output` | Formats responses in three modes via a `Formatter` interface. |
| `internal/ui` | Bubble Tea TUI helpers — spinners, progress watching, TTY/color detection. |

### Command Structure

Each command is a function returning `*cli.Command` (e.g. `geocodeCmd()`, `reverseCmd()`, `distanceCmd()`, `distanceMatrixCmd()`, `distanceJobsCmd()`, `listsCmd()`), all wired into the root command in `internal/cli/root.go` (`NewApp()`). Subcommands (e.g. `distance jobs create`, `lists upload`) are nested `*cli.Command` values returned by their own builder functions.

Available commands:
- `geocode` — forward geocode a single address or `--batch` a file (one address per line, up to 10,000)
- `reverse` — reverse geocode coordinates, single or batch
- `distance` — distance/travel time from an origin to one or more `-d` destinations
- `distance-matrix` — full N×M matrix
- `distance jobs` — async distance jobs: `create` / `list` / `status` / `download` / `delete`
- `lists` — spreadsheet geocoding: `upload` / `list` / `status` / `download` / `delete`

The typical action flow: build an `*App` via `newApp(cmd, ...)` (resolves config, constructs the API client and formatter), build a typed request struct from flags/args, call the corresponding `app.client.<Method>(ctx, req)`, then hand the response to `app.formatter.<Format...>(resp)`. Long-running batch calls are wrapped in `ui.WithSpinner(...)`.

### Configuration (`internal/config`)

- `GEOCODIO_API_KEY` env var or `--api-key` flag (flag wins). Missing key → `MissingAPIKeyError`.
- Base URL defaults to `https://api.geocod.io/v2`; override with `--base-url` (e.g. Enterprise hosts).
- `--debug` enables request/response logging to stderr.

### API client (`internal/api`)

- Single, plain `net/http`-based `Client`; functional options: `WithHTTPClient`, `WithUserAgent`, `WithDebug`.
- `api_key` is sent as a **query parameter**; default timeout is 30s.
- Single queries use **GET**; batch queries use **POST** with a JSON body (matching the Geocodio API).
- User-Agent is set to `geocodio-cli/<version>`.
- Errors are mapped to typed errors in `errors.go`.

### Output modes (`internal/output`)

A `Formatter` interface with three implementations, selected by global flags:
- **human** (default) — styled terminal tables; honors `--no-color` / `NO_COLOR` and TTY detection (`internal/ui`).
- **JSON** (`--json`) — raw JSON, for scripting/parsing.
- **agent** (`--agent`) — clean markdown, intended for LLM consumption.

Global flags live on the root command: `--api-key`, `--base-url`, `--json`, `--agent`, `--no-color`, `--debug`.

## Testing Strategy

- Tests run against recorded HTTP fixtures using **go-vcr v4** (`gopkg.in/dnaeon/go-vcr.v4`). Cassettes are YAML files in `internal/api/testdata/`.
- Default mode is **replay-only**, so `make test` needs no API key. Tests skip with a hint if a cassette is missing.
- To record/refresh cassettes, run `make record-cassettes` (sets `VCR_MODE=record`) with a real `GEOCODIO_API_KEY`. The API key is redacted from saved cassettes (`redactHook`) and ignored when matching requests (`matcherIgnoringAPIKey`) — see `internal/api/testutil_test.go`.
- `integration_test.go` (repo root) is a high-level end-to-end test built from recorded real API responses.
- Assertions use `testify`. `scripts/smoke-test.sh` provides a quick manual smoke check.

## Version Information

`Version` is a package var in `internal/cli` (`root.go`), defaulting to `"dev"` and overridden at build time:

```
-ldflags "-X github.com/geocodio/geocodio-cli/internal/cli.Version=$(VERSION)"
```

## Agent Skill

`skills/geocodio/SKILL.md` is a Claude Code skill teaching agents how to drive the `geocodio` binary (preferring `--json` for parsing and `--agent` for presenting results). Keep it in sync when commands, flags, or output formats change.

## Adding a New Command

1. Add a `fooCmd() *cli.Command` builder in a new (or existing) file under `internal/cli`.
2. Implement its `Action` (signature `func(ctx context.Context, cmd *cli.Command) error`); start it by calling `newApp(cmd, ...)`.
3. Register the command in the `Commands` slice of `NewApp()` in `internal/cli/root.go` (or nest it under a parent command's `Commands`).
4. Add the endpoint method and request/response types in `internal/api`.
5. Add a `Format…` method to the `Formatter` interface (`internal/output/formatter.go`) and implement it in `human.go`, `json.go`, and `agent.go`.
6. Add tests with a recorded cassette in `internal/api/testdata/`; update `skills/geocodio/SKILL.md` if the command is user-facing.
