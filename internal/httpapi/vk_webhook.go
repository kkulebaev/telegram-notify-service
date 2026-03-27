package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type vkWebhookRequest struct {
	Type         string          `json:"type"`
	GroupID      int             `json:"group_id"`
	Secret       *string         `json:"secret"`
	Object       json.RawMessage `json:"object"`
	EventID      *string         `json:"event_id"`
	V            *string         `json:"v"`
	Confirmation *string         `json:"confirmation"`
}

func (h *Handler) vkWebhook(w http.ResponseWriter, r *http.Request) {
	var req vkWebhookRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("decode json: %w", err))
		return
	}

	if strings.TrimSpace(req.Type) == "confirmation" {
		token := strings.TrimSpace(h.cfg.VKConfirmationToken)
		if token == "" {
			writeError(w, http.StatusBadRequest, fmt.Errorf("vk confirmation token is not configured"))
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(token))
		return
	}

	if strings.TrimSpace(h.cfg.VKSecret) != "" {
		got := ""
		if req.Secret != nil {
			got = strings.TrimSpace(*req.Secret)
		}
		if got != strings.TrimSpace(h.cfg.VKSecret) {
			writeError(w, http.StatusUnauthorized, fmt.Errorf("invalid vk secret"))
			return
		}
	}

	text := fmt.Sprintf("VK webhook event: %s", strings.TrimSpace(req.Type))
	var timestamp *time.Time
	if strings.TrimSpace(req.Type) == "message_new" {
		if parsed, ok := parseVKMessageTextAndTime(req.Object); ok {
			if strings.TrimSpace(parsed.Text) != "" {
				text = parsed.Text
			}
			if !parsed.Time.IsZero() {
				ts := parsed.Time
				timestamp = &ts
			}
		}
	}

	payload := map[string]any{
		"text":      text,
		"title":     "VK: Новое сообщение в группе",
		"timestamp": timestamp,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("marshal notify request: %w", err))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/api/v1/notify", strings.TrimRight(h.notifyBaseURL, "/"))

	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("create internal request: %w", err))
		return
	}

	hreq.Header.Set("Content-Type", "application/json")
	hreq.Header.Set("X-Admin-Token", h.cfg.AdminToken)

	resp, err := h.httpClient.Do(hreq)
	if err != nil {
		log.Error().Err(err).Msg("failed to call internal notify")
		writeError(w, http.StatusBadGateway, fmt.Errorf("call internal notify: %w", err))
		return
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode/100 != 2 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		writeError(w, http.StatusBadGateway, fmt.Errorf("internal notify returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body))))
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
