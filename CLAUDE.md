# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A minimal Go library for working with the Telegram Bot API. The library provides a lightweight client for interacting with Telegram bots.

**Current Implementation:**
- `Client` struct in `client.go` - main API client
- Single method: `SendMessage` for sending text messages to chats
- Uses standard `net/http` client with context support
- Bot token-based authentication

**Module:** `github.com/Roma7-7-7/telegram`

## Development Commands

### Go Commands
```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run a single test
go test -v -run TestName ./path/to/package

# Format code
go fmt ./...

# Lint code (requires golangci-lint)
golangci-lint run

# Install dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## Architecture

### Client Design
- **Constructor:** `NewClient(httpClient *http.Client, botToken string) *Client`
  - Accepts an `*http.Client` for flexibility (custom timeouts, transport, etc.)
  - Requires Telegram bot token for API authentication

- **API Methods:** All methods follow this pattern:
  - Accept `context.Context` as first parameter for cancellation/timeout control
  - Return `error` for error handling
  - Use internal request structs with JSON tags for API payloads
  - Construct URLs as `https://api.telegram.org/bot{token}/{method}`
  - Set `Content-Type: application/json` header
  - Check for HTTP 200 status code

### Code Style
- Error wrapping with `fmt.Errorf` and `%w` verb for error chains
- Use `//nolint:errcheck` with justification when intentionally ignoring errors
- Type definitions grouped in parenthesized blocks
- Request/response structs are private (lowercase) and co-located with methods

### Adding New API Methods

When adding support for new Telegram Bot API methods:

1. Define a private request struct with JSON tags (e.g., `sendMessageRequest`)
2. Create a public method on `*Client` with signature: `MethodName(ctx context.Context, params...) error` or `MethodName(ctx context.Context, params...) (*Result, error)`
3. Follow the existing pattern in `SendMessage` (client.go:30-56):
   - Marshal request to JSON
   - Create HTTP request with context
   - Set `Content-Type: application/json` header
   - Execute request
   - Check status code
   - Parse response if needed

### Dependencies

The library intentionally has zero external dependencies (only Go standard library). Maintain this constraint when adding features.
