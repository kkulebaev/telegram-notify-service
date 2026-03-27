# telegram-notify-service

A tiny HTTP service that sends notifications to Telegram.

## Endpoints

- `GET /healthz` → `200 ok`
- `POST /notify` → sends a message to Telegram

## Environment variables

- `TELEGRAM_BOT_TOKEN` *(required)*
- `TELEGRAM_CHAT_ID` *(required)*
- `PORT` *(optional, default: 8080)*

## Request example

```bash
curl -X POST http://localhost:8080/notify \
  -H 'content-type: application/json' \
  -d '{
    "text": "Deploy failed on migrations",
    "title": "Prod deploy failed",
    "level": "error",
    "source": "payments-api",
    "links": [{"label":"Logs","url":"https://example.com/logs"}]
  }'
```
