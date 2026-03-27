package httpapi

import (
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
	cfg    config.Config
	sender telegramapi.Sender
}

func NewHandler(cfg config.Config) *Handler {
	client := &http.Client{Timeout: 10 * time.Second}

	return &Handler{
		cfg:    cfg,
		sender: telegram.NewSender(client, cfg.TelegramBotToken, cfg.TelegramChatID),
	}
}

func NewHandlerWithSender(cfg config.Config, sender telegramapi.Sender) *Handler {
	return &Handler{cfg: cfg, sender: sender}
}

func (h *Handler) Router() http.Handler {
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
		r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		r.Post("/notify", h.notify)
	})

	return r
}
