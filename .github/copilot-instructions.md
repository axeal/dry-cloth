# dry-cloth Copilot Instructions

## Project Overview

dry-cloth is a **DigitalOcean droplet cleanup utility** that automatically deletes droplets older than a specified age. This is a single-purpose CLI tool, not a web service or long-running daemon.

**Key Architecture:**
- CLI entry point: [cmd/dry-cloth/main.go](cmd/dry-cloth/main.go) - handles argument parsing via `urfave/cli/v2`
- Core logic: [pkg/drycloth/drycloth.go](pkg/drycloth/drycloth.go) - implements droplet listing and deletion via DigitalOcean API
- Single operation: lists all droplets → filters by age → optionally skips by tag → deletes matching droplets

## Project-Specific Conventions

### Go Module Structure
- Module path: `github.com/axeal/dry-cloth`
- Uses Go 1.25+
- Primary dependency: `github.com/digitalocean/godo` for DigitalOcean API interactions

### CLI Design Pattern
- All configuration via flags or environment variables (no config files)
- Access token: `--access-token` flag or `DIGITALOCEAN_ACCESS_TOKEN` env var (required)
- Preservation logic: droplets with `--preserve-tag` are never deleted
- Safety mechanism: `--dry-run` flag for preview without deletion
- Default retention: 14 days via `--max-age-days` (overridable)

### Code Organization Patterns
- Business logic in `pkg/drycloth/` (library-style package)
- Pagination handled explicitly per DigitalOcean API requirements (see `DropletList` function)
- Date parsing uses RFC3339 format: `2006-01-02T15:04:05Z`
- Simple helper functions (`contains`, `parseDate`) defined locally rather than external dependencies

## Developer Workflows

### Building
```bash
# Local build
go build -o dry-cloth cmd/dry-cloth/main.go

# Docker build (multi-stage with distroless final image)
docker build -t dry-cloth .
```

### Testing Locally
Always use `--dry-run` first to preview deletions:
```bash
export DIGITALOCEAN_ACCESS_TOKEN=your_token
go run cmd/dry-cloth/main.go --preserve-tag DoNotDelete --max-age-days 28 --dry-run
```

### Deployment Context
- Containerized for production (see [Dockerfile](Dockerfile))
- Distroless base image for minimal attack surface
- Runs as non-root user
- Typically scheduled via cron/K8s CronJob (not included in this repo)

## Critical Implementation Details

### Pagination Requirement
The DigitalOcean API requires manual pagination handling. The `DropletList` function demonstrates the pattern:
1. Start with empty `ListOptions`
2. Loop through pages using `resp.Links.CurrentPage()`
3. Accumulate results until `IsLastPage()` returns true

### Error Handling Philosophy
- API errors are logged but don't stop processing of other droplets
- Individual droplet parsing/deletion errors are logged and skipped
- Only fatal errors (e.g., initial API connection) return early

### Tag-Based Preservation
The `preserveTag` check happens **after** age filtering. A droplet must be both old AND not have the preservation tag to be deleted. If preservation tag is not specified, no droplets are protected.

## When Making Changes

- **Adding new flags:** Update both the CLI app in `main.go` and the `Run` function signature in `drycloth.go`
- **Modifying deletion logic:** Core logic is in the `Run` function's droplet iteration loop
- **Changing date handling:** All dates from DigitalOcean API use RFC3339 format
- **Docker updates:** Remember the multi-stage build pattern - builder stage uses alpine, final uses distroless
