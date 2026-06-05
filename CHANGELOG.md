# Changelog

All notable changes to the Geocodio CLI are documented in this file.

## v3.0.0 - 2026-06-05

### Breaking changes

- Migrated the default Geocodio API version from `v1.9` to **v2**. Requests now go to `https://api.geocod.io/v2/...`.
- Removed the top-level `input` object from `/geocode` and `/reverse` responses. The parsed address now lives in each result under `address_components`.
- Renamed keys inside `address_components` (and `address_components_secondary`):
  - `zip` → `postal_code`
  - `state` → `state_province`
  - `secondaryunit` → `unit_type`
  - `secondarynumber` → `unit_number`
