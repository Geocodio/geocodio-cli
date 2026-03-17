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

## Commands

### Geocode (address to coordinates)

```bash
# Single address
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC"

# With data appends (timezone, congressional district, census, etc.)
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --fields timezone,cd

# Limit results
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --limit 1

# Country hint (US or CA)
geocodio geocode "Ottawa, Ontario" --country CA

# Batch from file (one address per line, max 10,000)
geocodio geocode --batch addresses.txt

# With inline distance to destinations
geocodio geocode "Washington DC" -d "New York" -d "Boston" --distance-mode driving

# Show stable address key
geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --show-address-key
```

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

# Options: --mode (driving|straightline), --units (miles|km)
geocodio distance "Washington DC" "New York" --mode driving --units km
```

### Distance Matrix (many-to-many)

```bash
# Requires files for origins and destinations (one location per line)
geocodio distance-matrix --origins origins.txt --destinations destinations.txt --mode driving --units miles
```

### Async Distance Jobs (large calculations)

```bash
geocodio distance-jobs create --name "My Job" --origins origins.txt --destinations destinations.txt --watch
geocodio distance-jobs list
geocodio distance-jobs status JOB_ID --watch
geocodio distance-jobs download JOB_ID --output results.csv
geocodio distance-jobs delete JOB_ID
```

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

**Process a CSV with address columns:**
```bash
geocodio lists upload customers.csv --direction forward --format "{{B}}, {{C}}, {{D}}" --watch
```

**Get timezone and congressional district for coordinates:**
```bash
geocodio reverse "38.8976,-77.0365" --fields timezone,cd --agent
```

## Batch File Limits

- Batch geocode/reverse: max 10,000 items per request
- For larger files, use `lists upload` which handles up to 10,000,000+ rows asynchronously
