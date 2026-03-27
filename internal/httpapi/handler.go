package httpapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kkulebaev/telegram-notify-service/internal/config"
	"github.com/kkulebaev/telegram-notify-service/internal/telegram"
	"github.com/kkulebaev/telegram-notify-service/internal/telegramapi"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	cfg           config.Config
	sender        telegramapi.Sender
	httpClient    *http.Client
	notifyBaseURL string
}

func NewHandler(cfg config.Config) *Handler {
	client := &http.Client{Timeout: 10 * time.Second}

	return &Handler{
		cfg:           cfg,
		sender:        telegram.NewSender(client, cfg.TelegramBotToken, cfg.TelegramChatID),
		httpClient:    client,
		notifyBaseURL: fmt.Sprintf("http://127.0.0.1:%d", cfg.Port),
	}
}

func NewHandlerWithSender(cfg config.Config, sender telegramapi.Sender) *Handler {
	client := &http.Client{Timeout: 10 * time.Second}

	return &Handler{
		cfg:           cfg,
		sender:        sender,
		httpClient:    client,
		notifyBaseURL: fmt.Sprintf("http://127.0.0.1:%d", cfg.Port),
	}
}

func (h *Handler) Router() http.Handler {
	if h.notifyBaseURL == "" {
		h.notifyBaseURL = fmt.Sprintf("http://127.0.0.1:%d", h.cfg.Port)
	}
	if h.httpClient == nil {
		h.httpClient = http.DefaultClient
	}

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(middleware.Compress(5))

	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: &log.Logger, NoColor: true}))
	// auth must be registered before any routes (chi requirement)
	r.Use(authMiddleware(h.cfg.AdminToken))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/webhooks/vk", h.vkWebhook)
		r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		r.Post("/notify", h.notify)
	})

	return r
}
