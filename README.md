# OmniLib

OmniLib hosts the data models and service definitions for all Omnisciens services, defining the contracts between services and the data models for backend services and ArangoDB. Protobuf files are published to https://buf.build/omnsight/omnlib.

## Local Development

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

## Development Guide

Tagging is automatically done by the github action. Commit message including `#major`, `#minor`, `#patch`, or `#none` in the branches `main` and `pre-release` will bump the release and pre-release versions.
