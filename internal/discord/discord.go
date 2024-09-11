package discord

import (
	"bytes"

	service "github.com/stodis/stodis/internal/service"
)

type DiscordFileService struct {
	bot     *Bot
	channel string
}

func NewDiscordFileService(bot *Bot, channel string) *DiscordFileService {
	return &DiscordFileService{
		bot:     bot,
		channel: channel,
	}
}

func (s *DiscordFileService) UploadFile(file []byte, name string) (string, error) {
	data, err := s.bot.Session.ChannelFileSend(s.channel, name, bytes.NewReader(file))
	if err != nil {
		return "", err
	}

	return data.ID, nil
}

var _ service.FileService = &DiscordFileService{}
