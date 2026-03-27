package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var ErrBadRequest = errors.New("telegram bad request")

type Sender struct {
	client *http.Client
	token  string
	chatID string
}

func NewSender(client *http.Client, token string, chatID string) *Sender {
	return &Sender{client: client, token: token, chatID: chatID}
}

type sendMessageRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type sendMessageResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

func (s *Sender) SendHTML(ctx context.Context, htmlMessage string) error {
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", url.PathEscape(s.token))

	body, err := json.Marshal(sendMessageRequest{
		ChatID:    s.chatID,
		Text:      htmlMessage,
		ParseMode: "HTML",
	})
	if err != nil {
		return fmt.Errorf("marshal telegram request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create telegram request: %w", err)
	}
	// Telegram wants JSON
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("send telegram request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

	if resp.StatusCode >= 400 {
		var parsed sendMessageResponse
		_ = json.Unmarshal(respBody, &parsed)
		if parsed.Description != "" {
			return fmt.Errorf("%w: %s", ErrBadRequest, parsed.Description)
		}
		return fmt.Errorf("%w: http %d", ErrBadRequest, resp.StatusCode)
	}

	return nil
}
