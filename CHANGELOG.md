# Changelog

All notable changes to the Geocodio CLI are documented in this file.

## Unreleased

### Changed

- **`--mode` now defaults to `straightline` instead of `driving`** on `distance`, `distance-matrix`, and `distance-jobs create`. This matches the HTTP API's own default and halves the cost of a run that didn't set `--mode`: straightline is 1 credit per calculation, driving is 2. Pass `--mode driving` explicitly if you need road distance or travel time. Note that `--max-duration` and `--min-duration` only apply in `driving` mode, so those filters now need `--mode driving` alongside them.

### Added

- **Cost notice before large distance runs.** `distance-matrix` and `distance-jobs create` print the calculation count, the routing mode, and the estimated credits to stderr when a request would exceed 10,000 calculations, along with a pointer to [setting a usage limit](https://www.geocod.io/guides/set-a-usage-limit).
- **Billing counters in `--json` and `--agent` output.** Responses now carry the `X-BILLABLE-LOOKUPS-COUNT` and `X-BILLABLE-DISTANCE-CALCULATIONS` headers as a `billing` object (`--json`) or a `_Billed: ..._` line (`--agent`), so callers can reconcile actual spend against their estimate.
- **Cost guidance in the bundled agent skill.** `skills/geocodio/SKILL.md` now documents lookups and credits, distance cost by mode, the free tier, `--batch` vs `lists upload`, cheaper paths (`--skip-geocoding`, `stable_address_key`), and country handling, and links to <https://www.geocod.io/docs/for-agents.md>.

- **Radius limiting for distance calculations.** `distance` and `distance-matrix` now accept `--max-distance` (aliased `--radius`), `--min-distance`, `--max-duration`, `--min-duration`, `--max-results`, `--order-by`, and `--sort-order`, mapping to the corresponding parameters on the `/distance` and `/distance-matrix` endpoints. Duration filters apply in `driving` mode only.
- `--distance-radius` as a friendlier alias for `geocode`/`reverse`'s existing `--distance-max-distance`.

### Fixed

- The `--batch` size error now reports how many records the file has, explains that data appends count toward the 10,000-lookup cap, and points at `lists upload` for larger files.
- Documented the inline distance filters on `geocode` and `reverse` (`--distance-max-distance`, `--distance-min-distance`, `--distance-max-duration`, `--distance-min-duration`, `--distance-max-results`, `--distance-order-by`, `--distance-sort-order`). They were already implemented but missing from the README and the agent skill.

## v3.2.0 - 26-08-25

### Added

- Friendly labels for UK data appends in the default (human) output. `uk-westminster`, `uk-devolved`, and `uk-local` now render as "Westminster Constituency", "Devolved Legislature", and "Local Authority" with the district name, instead of the raw field key.

## v3.1.0 - 2026-07-17

### Added

- **United Kingdom support.** Pass `--country "United Kingdom"` as a country hint, and request UK data appends via `--fields`: `uk-westminster`, `uk-westminster-next`, `uk-devolved`, `uk-devolved-next`, `uk-local`, `uk-local-next` (the `-next` variants return upcoming boundary changes).
- `match_type`, `address_lines`, and `formatted_street` are now surfaced in geocode results.

### Changed

- Removed client-side country validation — country hints are passed through to the API as-is.

## v3.0.0 - 2026-06-05

### Breaking changes

- Migrated the default Geocodio API version from `v1.9` to **v2**. Requests now go to `https://api.geocod.io/v2/...`.
- Removed the top-level `input` object from `/geocode` and `/reverse` responses. The parsed address now lives in each result under `address_components`.
- Renamed keys inside `address_components` (and `address_components_secondary`):
  - `zip` → `postal_code`
  - `state` → `state_province`
  - `secondaryunit` → `unit_type`
  - `secondarynumber` → `unit_number`
