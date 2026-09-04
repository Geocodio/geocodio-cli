---
name: geocodio
description: Use the Geocodio CLI to geocode addresses, reverse geocode coordinates, calculate distances, and process spreadsheets. Invoke when the user asks to geocode, reverse geocode, look up an address or coordinates, calculate distances or travel times, process a CSV/spreadsheet of addresses, or anything related to Geocodio's API. Requires the `geocodio` binary installed and `GEOCODIO_API_KEY` set.
---

# Geocodio CLI

Use the `geocodio` command-line tool to interact with the Geocodio API.

## Prerequisites

- `geocodio` binary installed and on PATH
- `GEOCODIO_API_KEY` environment variable set

If either is missing, tell the user what they need before proceeding.

## Output Format

Always use `--json` when you need to parse or process results programmatically. Use `--agent` when presenting results directly to the user in conversation (produces clean markdown tables). Omit both for default human-readable terminal output.

Both modes report what the request was actually billed — a `billing` object in `--json` (`lookups`, `distance_calculations`) and a `_Billed: ..._` line in `--agent`. Use it to reconcile actual spend against your estimate.

## Cost and Credits

Every Geocodio request costs credits. Read this before running anything large — the user pays for it.

**Lookups.** 1 lookup = 1 credit. A request costs `records x (1 + number of appends)` lookups: each `--fields` category is a full extra lookup per record, so 1,000 addresses with `--fields timezone,cd` is 3,000 lookups, not 1,000. Zero-result lookups are not billed. On pay-as-you-go and Flex plans the first 2,500 lookups per day are free.

**Distance calculations.** A distance request costs geocoding lookups for any inputs given as addresses (coordinates cost nothing to geocode), plus `origins x destinations x mode multiplier` distance calculations, plus any appends:

| Mode | Cost per calculation |
|------|----------------------|
| `straightline` (default) | 1 credit |
| `driving` | 2 credits |

The CLI defaults to `straightline`, matching the API. Only pass `--mode driving` when the user actually needs road distance or travel time — it doubles the bill. A 500x500 matrix is 250,000 credits straightline and 500,000 driving.

**Before a large run:** estimate the lookups, tell the user what it will cost, and get confirmation. Point them at the usage limit setting so a mistake can't run away: https://www.geocod.io/guides/set-a-usage-limit

**Cheaper paths:**

- `--skip-geocoding` on `reverse` when you already have coordinates and only need appends.
- `--show-address-key` returns a `stable_address_key`; store it so the same address is never geocoded twice.
- Drop `--fields` categories the user didn't ask for. Each one is a full extra lookup per record.
- Prefer `lists upload` over `--batch` for anything large (see below).

**Country.** Always pass `--country` explicitly. When it's omitted the API falls back to the US, which silently produces wrong (and billed) results for non-US addresses. UK geocoding requires a Flex or Unlimited+UK plan.

Full agent-facing cost guidance: https://www.geocod.io/AGENTS.md

## Commands

### Geocode (address to coordinates)

```bash
# Single address
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC"

# With data appends (timezone, congressional district, census, etc.)
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --fields timezone,cd

# Limit results
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --limit 1

# Country hint (e.g. USA, Canada, United Kingdom)
geocodio geocode "Ottawa, Ontario" --country Canada

# UK address with UK-specific data appends
geocodio geocode "10 Downing St, London" --country "United Kingdom" --fields uk-westminster,uk-local

# Batch from file (one address per line, max 10,000)
geocodio geocode --batch addresses.txt

# With inline distance to destinations
geocodio geocode "Washington DC" -d "New York" -d "Boston" --distance-mode driving

# Inline distance with a radius limit (note the distance- prefix on these commands)
geocodio geocode "Washington DC" -d "New York" -d "Boston" --distance-radius 150

# Show stable address key
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --show-address-key
```

UK data append fields: `uk-westminster`, `uk-westminster-next`, `uk-devolved`, `uk-devolved-next`, `uk-local`, `uk-local-next` (the `-next` variants return upcoming boundary changes). Pass them via `--fields` like any other append.

### Reverse Geocode (coordinates to address)

```bash
# Single coordinate (lat,lng)
geocodio reverse "38.8976,-77.0365"

# Batch from file (one lat,lng per line, max 10,000)
geocodio reverse --batch coordinates.txt

# Skip geocoding, only get field appends
geocodio reverse "38.8976,-77.0365" --skip-geocoding --fields timezone

# With inline distance to destinations
geocodio reverse "38.8976,-77.0365" -d "New York" --distance-mode driving
```

### Distance (origin to destinations)

```bash
# Single destination
geocodio distance "Washington DC" "New York"

# Multiple destinations
geocodio distance "Washington DC" "New York" "Boston" "Philadelphia"

# Options: --mode (straightline|driving), --units (miles|km)
# straightline is the default and costs 1 credit per calculation; driving costs 2
geocodio distance "Washington DC" "New York" --mode driving --units km
```

**Radius limiting and filtering.** Use these to answer "which of these are within X of the origin?" rather than filtering the full result set yourself:

```bash
# Only destinations within 150 miles (--radius is an alias for --max-distance)
geocodio distance "Washington DC" "New York" "Boston" "Philadelphia" --radius 150

# A ring: between 50 and 150 miles out
geocodio distance "Washington DC" "New York" "Boston" --min-distance 50 --max-distance 150

# Within a 2-hour drive (duration filters need --mode driving)
geocodio distance "Washington DC" "New York" "Boston" --mode driving --max-duration 7200

# The 3 nearest, sorted by driving time
geocodio distance "Washington DC" "New York" "Boston" "Philadelphia" \
  --mode driving --max-results 3 --order-by duration --sort-order asc
```

| Flag | Description |
|------|-------------|
| `--max-distance` / `--radius` | Radius limit: only keep destinations within this distance (in `--units`) |
| `--min-distance` | Only keep destinations at least this far away (in `--units`) |
| `--max-duration` | Only keep destinations within this travel time in seconds (driving mode only) |
| `--min-duration` | Only keep destinations at least this many seconds away (driving mode only) |
| `--max-results` | Only keep the N nearest destinations per origin |
| `--order-by` | Sort destinations by `distance` (default) or `duration` |
| `--sort-order` | `asc` (default) or `desc` |

### Distance Matrix (many-to-many)

```bash
# Requires files for origins and destinations (one location per line)
geocodio distance-matrix --origins origins.txt --destinations destinations.txt --mode driving --units miles

# The same radius/duration/sorting flags apply, per origin
geocodio distance-matrix --origins customers.txt --destinations stores.txt --radius 10

# Nearest store to each customer by driving time
geocodio distance-matrix --origins customers.txt --destinations stores.txt \
  --mode driving --max-results 1 --order-by duration
```

A matrix costs `origins x destinations` calculations. Above 10,000 calculations the CLI prints the count, the mode, and the estimated credits to stderr before submitting — read it back to the user rather than swallowing it.

### Async Distance Jobs (large calculations)

```bash
geocodio distance-jobs create --name "My Job" --origins origins.txt --destinations destinations.txt --watch
geocodio distance-jobs list
geocodio distance-jobs status JOB_ID --watch
geocodio distance-jobs download JOB_ID --output results.csv
geocodio distance-jobs delete JOB_ID
```

Distance jobs are billed the same way as `distance-matrix`, and print the same cost notice above 10,000 calculations. `--mode` defaults to `straightline` here too.

### Spreadsheet Processing (async batch geocoding)

```bash
# Upload CSV/Excel -- format uses {{A}}, {{B}}, {{C}} for columns
geocodio lists upload data.csv --direction forward --format "{{A}}, {{B}}, {{C}}" --watch

# With data append fields
geocodio lists upload data.csv --direction forward --format "{{A}}" --fields timezone,cd

# Reverse geocoding a spreadsheet
geocodio lists upload coords.csv --direction reverse --format "{{A}}" --watch

# Manage uploads
geocodio lists list
geocodio lists status LIST_ID --watch
geocodio lists download LIST_ID --output geocoded.csv
geocodio lists delete LIST_ID
```

**Concurrency.** Geocodio limits how many spreadsheet jobs one account can process at the same time. The exact number depends on your plan, and uploads started from the dashboard count toward the same limit as CLI/API uploads. There is no error or rejection for this: a list over the limit just stays in `ENQUEUED` longer and starts automatically once a slot frees up. If a list sits in `ENQUEUED` for a while, that's expected -- don't retry the upload or treat it as a failure. See https://www.geocod.io/guides/why-spreadsheet-uploads-get-queued for your plan's actual limit before running multiple lists at once.

## Global Flags

These work with all commands:

| Flag | Description |
|------|-------------|
| `--json` | Raw JSON output (for parsing) |
| `--agent` | Markdown output (for conversation) |
| `--api-key` | Override API key |
| `--no-color` | Disable colors |
| `--debug` | Show HTTP request/response details |

## Common Workflows

**Geocode a list of addresses and get JSON results:**
```bash
geocodio geocode --batch addresses.txt --json
```

**Find distances from one address to several others:**
```bash
geocodio geocode "123 Main St, Springfield IL" -d "Chicago IL" -d "St Louis MO" --distance-mode driving --agent
```

**Find which locations are within a radius:**
```bash
geocodio distance "123 Main St, Springfield IL" "Chicago IL" "St Louis MO" "Indianapolis IN" --radius 100 --agent
```

**Process a CSV with address columns:**
```bash
geocodio lists upload customers.csv --direction forward --format "{{B}}, {{C}}, {{D}}" --watch
```

**Get timezone and congressional district for coordinates:**
```bash
geocodio reverse "38.8976,-77.0365" --fields timezone,cd --agent
```

## Batch vs Lists

- `--batch` caps at **10,000 lookups** per request, and appends count toward that cap: 10,000 addresses with one `--fields` category is 20,000 lookups, which is over the cap.
- **Prefer `lists upload` for anything larger, and for any spreadsheet.** It is asynchronous, handles up to 10M lookups, and needs no splitting into chunks.
- Don't hand-split a large file into 10,000-line batches when `lists upload` would do it in one call.
