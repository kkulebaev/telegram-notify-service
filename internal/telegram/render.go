package telegram

import (
	"fmt"
	"html"
	"strings"
	"time"
)

type RenderParams struct {
	Level     Level
	Title     string
	Text      string
	Source    *string
	Links     []Link
	Timestamp time.Time
}

func RenderMessage(p RenderParams) string {
	icon, label := levelBadge(p.Level)

	lines := make([]string, 0, 16)
	lines = append(lines, fmt.Sprintf("%s <b>%s</b>", icon, html.EscapeString(p.Title)))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("📝 %s", html.EscapeString(p.Text)))

	_ = label

	if p.Source != nil {
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("🔎 <i>source:</i> <code>%s</code>", html.EscapeString(*p.Source)))
	}

	if len(p.Links) > 0 {
		lines = append(lines, "")
		lines = append(lines, "🔗 <b>Links</b>")
		for _, l := range p.Links {
			lines = append(lines, fmt.Sprintf("• <a href=\"%s\">%s</a>", html.EscapeString(l.URL), html.EscapeString(l.Label)))
		}
	}

	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("🕒 <i>%s UTC</i>", p.Timestamp.UTC().Format("2006-01-02 15:04")))

	return strings.Join(lines, "\n")
}

func levelBadge(l Level) (icon string, label string) {
	switch l {
	case LevelSuccess:
		return "✅", "SUCCESS"
	case LevelWarning:
		return "⚠️", "WARNING"
	case LevelError:
		return "🚨", "ERROR"
	default:
		return "ℹ️", "INFO"
	}
}
