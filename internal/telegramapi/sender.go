package telegramapi

import "context"

type Sender interface {
	SendHTML(ctx context.Context, htmlMessage string) error
}
