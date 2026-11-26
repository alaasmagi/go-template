# go-template

A Go backend template that mirrors [alaasmagi-dotnet-template](https://github.com/alaasmagi/alaasmagi-dotnet-template).
Standard library only. Conventions and architecture rules are in [AGENTS.md](AGENTS.md).

## Configuration

Copy `.env.example` to `.env` and fill in the values. The `.env` file is searched upward from the working
directory; variables already set in the process are never overridden. There are no defaults in code; a
missing or invalid value fails at startup.

| Variable | Description |
| --- | --- |
| `APP_PORT` | HTTP port (positive integer) |
| `APP_SLUG` | Application slug |

## Commands

```bash
go build ./...
go vet ./...
go test ./...
go run ./cmd/server
```

## Endpoints

```text
GET /health    returns 200 "ok"
```
