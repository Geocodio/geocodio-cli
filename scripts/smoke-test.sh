#!/usr/bin/env bash
#
# Smoke test for the Geocodio CLI.
#
# Runs every CLI command against the live Geocodio API to verify end-to-end
# behavior. This complements the unit tests (which use recorded VCR cassettes)
# by catching issues that only show up with real API responses.
#
# Prerequisites:
#   - GEOCODIO_API_KEY must be set (a real key — this hits the live API)
#   - The binary at ./bin/geocodio (auto-built if missing)
#
# Usage:
#   ./scripts/smoke-test.sh              # uses ./bin/geocodio
#   ./scripts/smoke-test.sh ./my-binary  # uses a custom binary path
#
# What it tests:
#   - geocode:          single, batch, all flags (fields, limit, country, destinations, output formats)
#   - reverse:          single, batch, all flags (fields, skip-geocoding, destinations, output formats)
#   - distance:         single/multi destination, driving/straightline, miles/km, output formats
#   - distance-matrix:  file-based origins/destinations, all mode/unit combos
#   - distance-jobs:    list, create → status → delete lifecycle
#   - lists:            list, upload → status → download → delete lifecycle
#   - global flags:     --version, --help, --no-color
#   - error handling:   missing API key, invalid coordinate format
#
# Note: This script creates and then deletes distance jobs and spreadsheet
# uploads during testing. These are cleaned up automatically, but will briefly
# appear in your Geocodio dashboard.

set -euo pipefail

CLI="${1:-./bin/geocodio}"
PASS=0
FAIL=0
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

# ─── Helpers ──────────────────────────────────────────────

red()   { printf "\033[31m%s\033[0m\n" "$*"; }
green() { printf "\033[32m%s\033[0m\n" "$*"; }
bold()  { printf "\033[1m%s\033[0m\n" "$*"; }

# run_test "Test name" command arg1 arg2 ...
# Runs the command and prints PASS/FAIL. On failure, shows the command and output.
run_test() {
    local name="$1"
    shift
    printf "  %-55s " "$name"
    if output=$("$@" 2>&1); then
        green "PASS"
        PASS=$((PASS + 1))
    else
        red "FAIL"
        echo "    Command: $*"
        echo "    Output:  $output"
        FAIL=$((FAIL + 1))
    fi
}

# ─── Preflight checks ────────────────────────────────────

if [ -z "${GEOCODIO_API_KEY:-}" ]; then
    red "GEOCODIO_API_KEY is not set. Export it before running this script."
    exit 1
fi

if [ ! -x "$CLI" ]; then
    echo "Binary not found at $CLI — building..."
    make build
fi

# ─── Temp fixture files (cleaned up on exit) ─────────────
cat > "$TMPDIR/addresses.txt" <<EOF
1600 Pennsylvania Ave NW, Washington DC
1 Infinite Loop, Cupertino CA
350 Fifth Avenue, New York NY
EOF

cat > "$TMPDIR/coordinates.txt" <<EOF
38.8976,-77.0365
37.3318,-122.0312
40.7484,-73.9857
EOF

cat > "$TMPDIR/origins.txt" <<EOF
Washington DC
Boston MA
EOF

cat > "$TMPDIR/destinations.txt" <<EOF
New York NY
Philadelphia PA
EOF

cat > "$TMPDIR/data.csv" <<EOF
address,city,state
1600 Pennsylvania Ave NW,Washington,DC
350 Fifth Avenue,New York,NY
EOF

# ─── Geocode ───────────────────────────────────────────────
bold "Geocode"
run_test "Single address"                    "$CLI" geocode "1600 Pennsylvania Ave NW, Washington DC"
run_test "Single address (JSON)"             "$CLI" geocode "1600 Pennsylvania Ave NW, Washington DC" --json
run_test "Single address (agent/markdown)"   "$CLI" geocode "1600 Pennsylvania Ave NW, Washington DC" --agent
run_test "With fields"                       "$CLI" geocode "1600 Pennsylvania Ave NW, Washington DC" --fields timezone
run_test "With limit"                        "$CLI" geocode "1600 Pennsylvania Ave NW, Washington DC" --limit 1
run_test "With country hint"                 "$CLI" geocode "1600 Pennsylvania Ave NW, Washington DC" --country US
run_test "With stable address key"           "$CLI" geocode "1600 Pennsylvania Ave NW, Washington DC" --show-address-key
run_test "With destinations"                 "$CLI" geocode "1600 Pennsylvania Ave NW, Washington DC" -d "New York" -d "Boston"
run_test "With destinations (driving)"       "$CLI" geocode "1600 Pennsylvania Ave NW, Washington DC" -d "New York" --distance-mode driving
run_test "With destinations (straightline)"  "$CLI" geocode "1600 Pennsylvania Ave NW, Washington DC" -d "New York" --distance-mode straightline
run_test "With destinations (km)"            "$CLI" geocode "1600 Pennsylvania Ave NW, Washington DC" -d "New York" --distance-units km
run_test "Batch geocode"                     "$CLI" geocode --batch "$TMPDIR/addresses.txt"
run_test "Batch geocode (JSON)"              "$CLI" geocode --batch "$TMPDIR/addresses.txt" --json
echo

# ─── Reverse Geocode ──────────────────────────────────────
bold "Reverse Geocode"
run_test "Single coordinate"                 "$CLI" reverse "38.8976,-77.0365"
run_test "Single coordinate (JSON)"          "$CLI" reverse "38.8976,-77.0365" --json
run_test "Single coordinate (agent)"         "$CLI" reverse "38.8976,-77.0365" --agent
run_test "With fields"                       "$CLI" reverse "38.8976,-77.0365" --fields timezone,cd
run_test "With limit"                        "$CLI" reverse "38.8976,-77.0365" --limit 1
run_test "With stable address key"           "$CLI" reverse "38.8976,-77.0365" --show-address-key
run_test "Skip geocoding (fields only)"      "$CLI" reverse "38.8976,-77.0365" --skip-geocoding --fields timezone
run_test "With destinations"                 "$CLI" reverse "38.8976,-77.0365" -d "New York"
run_test "Batch reverse"                     "$CLI" reverse --batch "$TMPDIR/coordinates.txt"
run_test "Batch reverse (JSON)"              "$CLI" reverse --batch "$TMPDIR/coordinates.txt" --json
echo

# ─── Distance ─────────────────────────────────────────────
bold "Distance"
run_test "Single destination"                "$CLI" distance "Washington DC" "New York"
run_test "Multiple destinations"             "$CLI" distance "Washington DC" "New York" "Boston" "Philadelphia"
run_test "Driving mode"                      "$CLI" distance "Washington DC" "New York" --mode driving
run_test "Straightline mode"                 "$CLI" distance "Washington DC" "New York" --mode straightline
run_test "Kilometers"                        "$CLI" distance "Washington DC" "New York" --units km
run_test "JSON output"                       "$CLI" distance "Washington DC" "New York" --json
run_test "Agent output"                      "$CLI" distance "Washington DC" "New York" --agent
echo

# ─── Distance Matrix ──────────────────────────────────────
bold "Distance Matrix"
run_test "Basic matrix"                      "$CLI" distance-matrix --origins "$TMPDIR/origins.txt" --destinations "$TMPDIR/destinations.txt"
run_test "Matrix (straightline)"             "$CLI" distance-matrix --origins "$TMPDIR/origins.txt" --destinations "$TMPDIR/destinations.txt" --mode straightline
run_test "Matrix (km)"                       "$CLI" distance-matrix --origins "$TMPDIR/origins.txt" --destinations "$TMPDIR/destinations.txt" --units km
run_test "Matrix (JSON)"                     "$CLI" distance-matrix --origins "$TMPDIR/origins.txt" --destinations "$TMPDIR/destinations.txt" --json
echo

# ─── Distance Jobs ────────────────────────────────────────
bold "Distance Jobs"
run_test "List jobs"                         "$CLI" distance-jobs list
run_test "List jobs (JSON)"                  "$CLI" distance-jobs list --json

# Lifecycle test: create a job, check its status, then clean it up.
# Uses --json output to parse the job ID for subsequent commands.
printf "  %-55s " "Create → status → delete lifecycle"
JOB_OUTPUT=$("$CLI" distance-jobs create \
    --name "smoke-test-$(date +%s)" \
    --origins "$TMPDIR/origins.txt" \
    --destinations "$TMPDIR/destinations.txt" \
    --json 2>&1) || true

JOB_ID=$(echo "$JOB_OUTPUT" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
if [ -n "$JOB_ID" ]; then
    STATUS_OK=true
    "$CLI" distance-jobs status "$JOB_ID" >/dev/null 2>&1 || STATUS_OK=false
    DELETE_OK=true
    "$CLI" distance-jobs delete "$JOB_ID" >/dev/null 2>&1 || DELETE_OK=false
    if $STATUS_OK && $DELETE_OK; then
        green "PASS"
        PASS=$((PASS + 1))
    else
        red "FAIL (status=$STATUS_OK, delete=$DELETE_OK)"
        FAIL=$((FAIL + 1))
    fi
else
    red "FAIL (could not parse job ID)"
    echo "    Output: $JOB_OUTPUT"
    FAIL=$((FAIL + 1))
fi
echo

# ─── Lists (Spreadsheets) ────────────────────────────────
bold "Lists"
run_test "List spreadsheets"                 "$CLI" lists list
run_test "List spreadsheets (JSON)"          "$CLI" lists list --json

# Lifecycle test: upload a CSV, check its status, attempt download, then clean up.
# Download may fail if the file is still processing — that's expected and ignored.
printf "  %-55s " "Upload → status → download → delete lifecycle"
LIST_OUTPUT=$("$CLI" lists upload "$TMPDIR/data.csv" \
    --direction forward \
    --format "{{A}}, {{B}}, {{C}}" \
    --json 2>&1) || true

LIST_ID=$(echo "$LIST_OUTPUT" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
if [ -n "$LIST_ID" ]; then
    STATUS_OK=true
    "$CLI" lists status "$LIST_ID" >/dev/null 2>&1 || STATUS_OK=false

    # Try download (may fail if still processing — that's okay)
    "$CLI" lists download "$LIST_ID" >/dev/null 2>&1 || true

    DELETE_OK=true
    "$CLI" lists delete "$LIST_ID" >/dev/null 2>&1 || DELETE_OK=false
    if $STATUS_OK && $DELETE_OK; then
        green "PASS"
        PASS=$((PASS + 1))
    else
        red "FAIL (status=$STATUS_OK, delete=$DELETE_OK)"
        FAIL=$((FAIL + 1))
    fi
else
    red "FAIL (could not parse list ID)"
    echo "    Output: $LIST_OUTPUT"
    FAIL=$((FAIL + 1))
fi
echo

# ─── Global flags ─────────────────────────────────────────
bold "Global Flags"
run_test "--version"                         "$CLI" --version
run_test "--help"                            "$CLI" --help
run_test "--no-color"                        "$CLI" geocode "1600 Pennsylvania Ave NW, Washington DC" --no-color
echo

# ─── Error handling ───────────────────────────────────────
bold "Error Handling"
printf "  %-55s " "Missing API key"
if "$CLI" geocode "test" --api-key "" 2>&1 | grep -qi "api key"; then
    green "PASS"
    PASS=$((PASS + 1))
else
    red "FAIL"
    FAIL=$((FAIL + 1))
fi

printf "  %-55s " "Invalid coordinate format"
if "$CLI" reverse "not-a-coordinate" 2>&1 | grep -qi "invalid\|error\|format"; then
    green "PASS"
    PASS=$((PASS + 1))
else
    red "FAIL"
    FAIL=$((FAIL + 1))
fi
echo

# ─── Summary ──────────────────────────────────────────────
echo
bold "Results: $PASS passed, $FAIL failed, $((PASS + FAIL)) total"
[ "$FAIL" -eq 0 ] && green "All tests passed!" || red "Some tests failed."
exit "$FAIL"
