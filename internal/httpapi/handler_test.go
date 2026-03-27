package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkulebaev/telegram-notify-service/internal/config"
)

type noopSender struct{}

func (noopSender) SendHTML(_ context.Context, _ string) error { return nil }

func TestHealthzRequiresAuth(t *testing.T) {
	cfg := config.Config{
		Port:             8080,
		TelegramBotToken: "x",
		TelegramChatID:   "1",
		AdminToken:       "secret",
	}

	h := NewHandlerWithSender(cfg, noopSender{})
	srv := httptest.NewServer(h.Router())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/healthz")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}
