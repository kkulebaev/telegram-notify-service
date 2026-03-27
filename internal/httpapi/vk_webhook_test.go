package httpapi

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/kkulebaev/telegram-notify-service/internal/config"
)

type spySender struct {
	mu    sync.Mutex
	calls []string
}

func (s *spySender) SendHTML(_ context.Context, msg string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = append(s.calls, msg)
	return nil
}

func (s *spySender) callCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.calls)
}

func TestVKWebhookEventReturnsOkAndSendsNotify(t *testing.T) {
	cfg := config.Config{
		Port:                8080,
		TelegramBotToken:    "x",
		TelegramChatID:      "1",
		AdminToken:          "secret",
		VKConfirmationToken: "confirm-token",
		VKSecret:            "vk-secret",
	}

	sender := &spySender{}
	h := NewHandlerWithSenderAndOptions(cfg, sender)
	srv := httptest.NewServer(h.Router())
	defer srv.Close()

	h.notifyBaseURL = srv.URL

	body := []byte(`{"type":"message_new","group_id":1,"secret":"vk-secret","object":{}}`)
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/webhooks/vk", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Admin-Token", "secret")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	got, _ := io.ReadAll(resp.Body)
	if string(got) != "ok" {
		t.Fatalf("expected body ok, got %q", string(got))
	}

	if sender.callCount() != 1 {
		t.Fatalf("expected sender calls 1, got %d", sender.callCount())
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestVKWebhookConfirmationReturnsToken(t *testing.T) {
	cfg := config.Config{
		Port:                8080,
		TelegramBotToken:    "x",
		TelegramChatID:      "1",
		AdminToken:          "secret",
		VKConfirmationToken: "confirm-token",
	}

	sender := &spySender{}
	h := NewHandlerWithSenderAndOptions(cfg, sender)
	srv := httptest.NewServer(h.Router())
	defer srv.Close()

	h.notifyBaseURL = srv.URL

	body := []byte(`{"type":"confirmation","group_id":1,"object":{}}`)
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/webhooks/vk", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Admin-Token", "secret")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	got, _ := io.ReadAll(resp.Body)
	if string(got) != "confirm-token" {
		t.Fatalf("expected body confirm-token, got %q", string(got))
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if sender.callCount() != 0 {
		t.Fatalf("expected sender calls 0, got %d", sender.callCount())
	}
}

func TestVKWebhookRejectsInvalidSecret(t *testing.T) {
	cfg := config.Config{
		Port:                8080,
		TelegramBotToken:    "x",
		TelegramChatID:      "1",
		AdminToken:          "secret",
		VKConfirmationToken: "confirm-token",
		VKSecret:            "vk-secret",
	}

	sender := &spySender{}
	h := NewHandlerWithSenderAndOptions(cfg, sender)
	srv := httptest.NewServer(h.Router())
	defer srv.Close()

	h.notifyBaseURL = srv.URL

	body := []byte(`{"type":"message_new","group_id":1,"secret":"wrong","object":{}}`)
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/webhooks/vk", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Admin-Token", "secret")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}

	if sender.callCount() != 0 {
		t.Fatalf("expected sender calls 0, got %d", sender.callCount())
	}
}
