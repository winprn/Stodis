package discord

import (
	"bytes"

	"github.com/stodis/stodis/internal/handler"
)

type DiscordHandler struct {
	bot     *Bot
	channel string
}

func NewDiscordHandler(bot *Bot, channel string) *DiscordHandler {
	return &DiscordHandler{
		bot:     bot,
		channel: channel,
	}
}

func (s *DiscordHandler) UploadFile(file []byte, name string) (string, error) {
	data, err := s.bot.Session.ChannelFileSend(s.channel, name, bytes.NewReader(file))
	if err != nil {
		return "", err
	}

	return data.ID, nil
}

var _ handler.StorageHandler = &DiscordHandler{}
