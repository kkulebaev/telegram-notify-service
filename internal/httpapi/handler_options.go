package httpapi

import (
	"net/http"

	"github.com/kkulebaev/telegram-notify-service/internal/config"
	"github.com/kkulebaev/telegram-notify-service/internal/telegramapi"
)

type HandlerOption func(*Handler)

func WithHTTPClient(c *http.Client) HandlerOption {
	return func(h *Handler) {
		h.httpClient = c
	}
}

func WithNotifyBaseURL(url string) HandlerOption {
	return func(h *Handler) {
		h.notifyBaseURL = url
	}
}

func NewHandlerWithSenderAndOptions(cfg config.Config, sender telegramapi.Sender, opts ...HandlerOption) *Handler {
	h := NewHandlerWithSender(cfg, sender)
	for _, opt := range opts {
		opt(h)
	}
	return h
}
