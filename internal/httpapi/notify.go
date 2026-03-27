package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kkulebaev/telegram-notify-service/internal/telegram"
	"github.com/rs/zerolog/log"
)

type notifyRequest struct {
	Text      string          `json:"text" validate:"required"`
	Title     *string         `json:"title" validate:"omitempty"`
	Level     *telegram.Level `json:"level" validate:"omitempty"`
	Source    *string         `json:"source" validate:"omitempty"`
	Links     []telegram.Link `json:"links" validate:"omitempty,dive"`
	Timestamp *time.Time      `json:"timestamp" validate:"omitempty"`
}

type notifyResponse struct {
	Ok bool `json:"ok"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func (h *Handler) notify(w http.ResponseWriter, r *http.Request) {
	var req notifyRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("decode json: %w", err))
		return
	}

	if err := validate.Struct(req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("validate body: %w", err))
		return
	}

	level := telegram.LevelInfo
	if req.Level != nil {
		level = *req.Level
	}

	title := "Notification"
	if req.Title != nil {
		title = strings.TrimSpace(*req.Title)
		if title == "" {
			title = "Notification"
		}
	}

	text := strings.TrimSpace(req.Text)
	timestamp := time.Now().UTC()
	if req.Timestamp != nil {
		timestamp = req.Timestamp.UTC()
	}

	msg := telegram.RenderMessage(telegram.RenderParams{
		Level:     level,
		Title:     title,
		Text:      text,
		Source:    optionalTrim(req.Source),
		Links:     req.Links,
		Timestamp: timestamp,
	})

	if err := h.sender.SendHTML(r.Context(), msg); err != nil {
		if errors.Is(err, telegram.ErrBadRequest) {
			log.Warn().Err(err).Msg("telegram rejected the request")
			writeError(w, http.StatusBadRequest, err)
			return
		}

		log.Error().Err(err).Msg("failed to send telegram message")
		writeError(w, http.StatusBadGateway, err)
		return
	}

	writeJSON(w, http.StatusOK, notifyResponse{Ok: true})
}

func optionalTrim(v *string) *string {
	if v == nil {
		return nil
	}

	s := strings.TrimSpace(*v)
	if s == "" {
		return nil
	}

	return &s
}
