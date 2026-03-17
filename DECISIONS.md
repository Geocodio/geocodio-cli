# Internal Design Decisions

This document tracks design decisions and open discussions for the CLI. **Not intended for public distribution.**

> **Note:** If this repo goes public, this file and the git history (including earlier commits where these notes were in the README) will be visible. Before open-sourcing, decide whether to squash commit history or accept the transparency.

## API Features Not Included

The Geocodio API supports some features that are intentionally not exposed in this CLI.

| Feature | API Support | CLI Status | Rationale |
|---------|------------|------------|-----------|
| Component address fields (`street`, `city`, `state`, etc.) | `GET /geocode` accepts individual address components as separate query params | Not included, except for `country` as a hint flag | The single address string covers the vast majority of use cases. Users with structured data can concatenate fields. The CLI still exposes `--country` because it is a useful disambiguation hint without requiring fully structured address input. Can be added later if there's demand. |
| `format=simple` | Single geocode/reverse requests support `format=simple` for simplified JSON output | Not included | The CLI already provides multiple output layers: human-readable (default), `--json` (full response), and `--agent` (markdown). A simplified JSON format would overlap with existing options. |
| Distance filtering params (`max_results`, `max_distance`, `max_duration`, `min_distance`, `min_duration`, `order_by`, `sort_order`) on standalone `GET /distance` | Supported on the distance endpoint | Not included | Overkill for a CLI. Users needing that level of filtering are likely using the API directly in code. |
| Distance filtering params on `POST /distance-matrix` | Supported on the distance-matrix endpoint | Not included | Same rationale as standalone distance. The filtering params exist on the `geocode --destinations` and `reverse --destinations` inline distance features where they make more sense for interactive use. |
| Distance filtering params, `fields`, and `callback_url` on `POST /distance-jobs` | Supported on the distance-jobs create endpoint | Not included | Distance jobs are for large async calculations. Users creating jobs of that scale are likely using the API directly. The CLI supports create/status/download/delete for basic job management. |
| Stable address keys shown by default | The API returns a `stable_address_key` on most geocode results | Hidden by default, opt-in with `--show-address-key` | Most users won't know what a stable address key is and showing it on every request would be confusing. Available via flag when needed. Always present in `--json` output. |

## Migration from v1.x: `lists` Naming

The v1.x CLI only had spreadsheet commands, so top-level names like `create` and `status` made sense. The v2.x CLI groups them under `lists` to match the actual API endpoint name (`/lists`), making it consistent with the Geocodio API docs. The new CLI also has more commands (geocode, reverse, distance), so grouping avoids ambiguity. This only affects existing spreadsheet users, but all of them are affected.

### Option: backward-compatible aliases

We could register the old command names (`create`, `status`, `download`, `remove`, `list`) as hidden aliases that map to the new `lists` subcommands. Existing scripts would keep working, and new users would only see the new names in help output. However, `status` and `download` are now ambiguous — they exist under both `lists` and `distance-jobs`. The aliases would point to `lists`, but this could be confusing if a user discovers them and tries to use them for distance jobs.
