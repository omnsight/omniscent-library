# OmniLib

Omnilib hosts the common data model for Omniscense micro-services.

## Local Development

Tag is injested by a github action. Commit message including `#major`, `#minor`, `#patch`, or `#none` will bump the release and pre-release versions.

### Dependencies

Buf build:

```bash
# Generate Go bindings for GeoVision services
buf registry login buf.build

buf dep update

buf format -w
buf lint

buf generate

buf push

go mod tidy
```
