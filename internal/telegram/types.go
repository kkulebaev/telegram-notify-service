package telegram

type Level string

const (
	LevelInfo    Level = "info"
	LevelWarning Level = "warning"
	LevelError   Level = "error"
	LevelSuccess Level = "success"
)

type Link struct {
	Label string `json:"label" validate:"required"`
	URL   string `json:"url" validate:"required,url"`
}
