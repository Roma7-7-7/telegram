package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrForbidden       = errors.New("forbidden")
	ErrTooManyRequests = errors.New("too many requests")
)

type (
	Context interface {
		context.Context
		ChatID() (string, bool)
	}

	MiddlewareFunc func(next Handler) Handler

	Handler func(ctx Context) error

	Option func(o *options)

	Client struct {
		httpClient *http.Client
		botToken   string
		options    *options
	}

	middlewareContext struct {
		context.Context
		chatID string
	}

	options struct {
		middlewares []MiddlewareFunc
	}

	sendMessageRequest struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}
)

func (c *middlewareContext) ChatID() (string, bool) {
	return c.chatID, c.chatID != ""
}

func WithMiddlewares(middlewares ...MiddlewareFunc) Option {
	return func(o *options) {
		o.middlewares = append(o.middlewares, middlewares...)
	}
}

func NewClient(httpClient *http.Client, botToken string, with ...Option) *Client {
	opts := &options{}
	for _, option := range with {
		option(opts)
	}

	return &Client{
		httpClient: httpClient,
		botToken:   botToken,
		options:    opts,
	}
}

func (c *Client) SendMessage(ctx context.Context, chatID, text string) error {
	return c.doWithMiddlewares(&middlewareContext{Context: ctx, chatID: chatID}, func(ctx Context) error {
		return c.sendMessage(ctx, chatID, text)
	})
}

func (c *Client) sendMessage(ctx Context, chatID, text string) error {
	body, err := json.Marshal(sendMessageRequest{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.botToken), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck // ignore

	//nolint:mnd // 1xx and 2xx status codes are ok
	if resp.StatusCode < 300 {
		return nil
	}

	switch resp.StatusCode {
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusTooManyRequests:
		return ErrTooManyRequests
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

func (c *Client) doWithMiddlewares(ctx Context, h Handler) error {
	for i := len(c.options.middlewares) - 1; i >= 0; i-- {
		h = c.options.middlewares[i](h)
	}
	return h(ctx)
}
