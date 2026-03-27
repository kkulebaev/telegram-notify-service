package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkulebaev/telegram-notify-service/internal/config"
	"github.com/kkulebaev/telegram-notify-service/internal/telegram"
)

type captureSender struct {
	called  bool
	message string
	err     error
}

func (s *captureSender) SendHTML(_ context.Context, htmlMessage string) error {
	s.called = true
	s.message = htmlMessage
	return s.err
}

func TestNotifyRequiresAuth(t *testing.T) {
	cfg := config.Config{Port: 8080, TelegramBotToken: "x", TelegramChatID: "1", AdminToken: "secret"}
	sender := &captureSender{}

	h := NewHandlerWithSender(cfg, sender)
	srv := httptest.NewServer(h.Router())
	defer srv.Close()

	body, _ := json.Marshal(map[string]any{"text": "hello"})
	resp, err := http.Post(srv.URL+"/notify", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
	if sender.called {
		t.Fatalf("sender must not be called when unauthorized")
	}
}

func TestNotifySendsMessage(t *testing.T) {
	cfg := config.Config{Port: 8080, TelegramBotToken: "x", TelegramChatID: "1", AdminToken: "secret"}
	sender := &captureSender{}

	h := NewHandlerWithSender(cfg, sender)
	srv := httptest.NewServer(h.Router())
	defer srv.Close()

	payload := map[string]any{
		"text":   "Deploy failed",
		"title":  "Prod",
		"level":  string(telegram.LevelError),
		"source": "payments-api",
		"links": []map[string]any{{
			"label": "Logs",
			"url":   "https://example.com/logs",
		}},
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/notify", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer secret")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if !sender.called {
		t.Fatalf("sender must be called")
	}
	if sender.message == "" {
		t.Fatalf("expected non-empty message")
	}
	if got, want := sender.message, "🚨 <b>Prod</b>"; !bytes.Contains([]byte(got), []byte(want)) {
		t.Fatalf("expected message to contain %q, got: %s", want, got)
	}
}

func TestNotifyRejectsUnknownFields(t *testing.T) {
	cfg := config.Config{Port: 8080, TelegramBotToken: "x", TelegramChatID: "1", AdminToken: "secret"}
	sender := &captureSender{}

	h := NewHandlerWithSender(cfg, sender)
	srv := httptest.NewServer(h.Router())
	defer srv.Close()

	body, _ := json.Marshal(map[string]any{"text": "hello", "nope": true})
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/notify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer secret")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
	if sender.called {
		t.Fatalf("sender must not be called on bad request")
	}
}
