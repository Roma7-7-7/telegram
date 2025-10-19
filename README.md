# Telegram Bot API Client

A minimal Go library for working with the Telegram Bot API.

## Features

- Lightweight and zero external dependencies (only Go standard library)
- Context-aware API methods for cancellation and timeout control
- Flexible HTTP client configuration
- Clean and simple API design

## Installation

```bash
go get github.com/Roma7-7-7/telegram
```

## Usage

### Creating a Client

```go
package main

import (
    "net/http"
    "time"

    "github.com/Roma7-7-7/telegram"
)

func main() {
    // Create an HTTP client with custom timeout
    httpClient := &http.Client{
        Timeout: 10 * time.Second,
    }

    // Initialize the Telegram client with your bot token
    client := telegram.NewClient(httpClient, "YOUR_BOT_TOKEN")
}
```

### Sending Messages

```go
package main

import (
    "context"
    "log"
    "net/http"
    "time"

    "github.com/Roma7-7-7/telegram"
)

func main() {
    httpClient := &http.Client{
        Timeout: 10 * time.Second,
    }

    client := telegram.NewClient(httpClient, "YOUR_BOT_TOKEN")

    // Send a message to a chat
    ctx := context.Background()
    err := client.SendMessage(ctx, "CHAT_ID", "Hello, World!")
    if err != nil {
        log.Fatalf("Failed to send message: %v", err)
    }
}
```

### Using Context for Timeout Control

```go
// Create a context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := client.SendMessage(ctx, "CHAT_ID", "Message with timeout")
if err != nil {
    log.Printf("Error: %v", err)
}
```

## API Methods

Currently supported Telegram Bot API methods:

- `SendMessage(ctx context.Context, chatID, text string) error` - Send text messages to a chat

## Contributing

This library aims to remain minimal and dependency-free. When adding new features:

1. Maintain zero external dependencies
2. Follow existing patterns and conventions
3. Include context support for all API methods
4. Use proper error wrapping with `fmt.Errorf` and `%w`

## License

MIT License - see [LICENSE](LICENSE) file for details.
