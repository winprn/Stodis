package discord

import (
	"bytes"
	"log"

	service "github.com/stodis/stodis/internal/service"
)

type DiscordFileService struct {
	bot     *Bot
	channel string
}

func NewDiscordFileService(bot *Bot, channel string) *DiscordFileService {
	bot.Session.Debug = true
	return &DiscordFileService{
		bot:     bot,
		channel: channel,
	}
}

func (s *DiscordFileService) UploadFile(file []byte, name string) (string, error) {
	log.Printf("Uploading file %s to Discord\n", name)
	data, err := s.bot.Session.ChannelFileSend(s.channel, name, bytes.NewReader(file))
	if err != nil {
		return "", err
	}

	return data.ID, nil
}

var _ service.FileService = &DiscordFileService{}
