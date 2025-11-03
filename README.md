# OmniLib

OmniLib hosts the data models and service definitions for all Omnisciens services, defining the contracts between services and the data models for backend services and ArangoDB.

## Directory Structure

```
gen/
├── go/ # go bindings
├── http/ # open api definitions
prots/
├── geovision/
├── google
├── model/
```

## Local Development

```bash
# Generate Go bindings for GeoVision services
make gen_geovision
make gen_geovision_openapi
go mod tidy
```

```bash
# Set up pre-commit hooks
git config core.hooksPath .githooks 
```

## Documentation Guidelines

Each readme file should cover ONLY these aspects:
- Project overview stating what the project is about
- Folder structure that covers what each main folder/file does
- Infrastructure overview that covers high-level design
- How to run locally

Each document should NOT cover any design/implementation details. The code should be self-explanatory.
