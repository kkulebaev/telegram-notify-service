package httpapi

import (
	"encoding/json"
	"time"
)

type vkParsedMessage struct {
	Text string
	Time time.Time
}

type vkMessageNewObject struct {
	Message *struct {
		Date *int64  `json:"date"`
		Text *string `json:"text"`
	} `json:"message"`
}

func parseVKMessageTextAndTime(raw json.RawMessage) (vkParsedMessage, bool) {
	var obj vkMessageNewObject
	if err := json.Unmarshal(raw, &obj); err != nil {
		return vkParsedMessage{}, false
	}
	if obj.Message == nil {
		return vkParsedMessage{}, false
	}

	res := vkParsedMessage{}
	if obj.Message.Text != nil {
		res.Text = *obj.Message.Text
	}
	if obj.Message.Date != nil {
		res.Time = time.Unix(*obj.Message.Date, 0).UTC()
	}

	return res, true
}
