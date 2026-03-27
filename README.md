<p align="center">
  <img src="./assets/gopher-telegram.svg" alt="telegram-notify-service" width="900" />
</p>

<p align="center">
  <a href="https://github.com/kkulebaev/telegram-notify-service/actions"><img alt="CI" src="https://img.shields.io/github/actions/workflow/status/kkulebaev/telegram-notify-service/ci.yml?branch=main" /></a>
  <a href="https://github.com/kkulebaev/telegram-notify-service/blob/main/LICENSE"><img alt="License" src="https://img.shields.io/github/license/kkulebaev/telegram-notify-service" /></a>
  <a href="https://github.com/kkulebaev/telegram-notify-service"><img alt="Go" src="https://img.shields.io/badge/Go-1.22-00ADD8?logo=go&logoColor=white" /></a>
  <a href="https://github.com/kkulebaev/telegram-notify-service"><img alt="Docker" src="https://img.shields.io/badge/Docker-ready-2496ED?logo=docker&logoColor=white" /></a>
</p>

# telegram-notify-service

Tiny HTTP service (Go) to send notifications to Telegram.

**Production URL:** https://telegram-notify-service-production.up.railway.app

## Features

- `POST /notify` — sends a formatted Telegram message (HTML + emoji)
- `GET /healthz` — healthcheck endpoint
- `ADMIN_TOKEN` protection (Bearer token / `X-Admin-Token`)
- Docker-friendly (distroless runtime image)

## Endpoints

### `GET /healthz`

Returns `200 ok`.

### `POST /notify`

Sends message to Telegram using `parse_mode=HTML`.

Request body:

```json
{
  "text": "Deploy failed on migrations",
  "title": "Prod deploy failed",
  "level": "error",
  "source": "payments-api",
  "links": [{ "label": "Logs", "url": "https://example.com/logs" }]
}
```

Supported fields:

- `text` *(string, required)*
- `title` *(string, optional)*
- `level` *("info" | "warning" | "error" | "success", optional; default: "info")*
- `source` *(string, optional)*
- `links` *(array, optional)* — `{ label, url }`
- `timestamp` *(ISO-8601 string, optional)*

## Configuration

Environment variables:

- `TELEGRAM_BOT_TOKEN` *(required)*
- `TELEGRAM_CHAT_ID` *(required)*
- `ADMIN_TOKEN` *(required)*
- `PORT` *(optional, default: 8080)*

See: [`.env.example`](./.env.example)

## Run locally

```bash
export TELEGRAM_BOT_TOKEN="..."
export TELEGRAM_CHAT_ID="..."
export ADMIN_TOKEN="..."
export PORT=8080

go run ./cmd/server
```

Test request:

```bash
curl -X POST https://telegram-notify-service-production.up.railway.app/notify \
  -H 'content-type: application/json' \
  -H 'authorization: Bearer <ADMIN_TOKEN>' \
  -d '{
    "text": "Deploy failed on migrations",
    "title": "Prod deploy failed",
    "level": "error",
    "source": "payments-api",
    "links": [{"label":"Logs","url":"https://example.com/logs"}]
  }'
```

## Run with Docker

```bash
docker build -t telegram-notify-service .

docker run --rm -p 8080:8080 \
  -e TELEGRAM_BOT_TOKEN="..." \
  -e TELEGRAM_CHAT_ID="..." \
  -e ADMIN_TOKEN="..." \
  telegram-notify-service
```

## Notes

- Telegram message delivery is synchronous: `POST /notify` returns error if Telegram API responds with an error.
- For production, consider adding rate limiting and deduplication depending on your alert sources.
