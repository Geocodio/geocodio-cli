# Changelog

All notable changes to the Geocodio CLI are documented in this file.

## Unreleased

### Added

- **Radius limiting for distance calculations.** `distance` and `distance-matrix` now accept `--max-distance` (aliased `--radius`), `--min-distance`, `--max-duration`, `--min-duration`, `--max-results`, `--order-by`, and `--sort-order`, mapping to the corresponding parameters on the `/distance` and `/distance-matrix` endpoints. Duration filters apply in `driving` mode only.
- `--distance-radius` as a friendlier alias for `geocode`/`reverse`'s existing `--distance-max-distance`.

### Fixed

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
