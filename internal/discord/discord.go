package discord

import (
	"bytes"
	"log"
	"time"

	service "github.com/stodis/stodis/internal/service"
)

type DiscordFileService struct {
	bot     *Bot
	channel string
}

func NewDiscordFileService(bot *Bot, channel string) *DiscordFileService {
	// bot.Session.Debug = true
	return &DiscordFileService{
		bot:     bot,
		channel: channel,
	}
}

func (s *DiscordFileService) UploadFile(file []byte, name string) (string, error) {
	log.Printf("%s is uploading file %s to Discord\n", s.bot.botID, name)
	startTime := time.Now()
	data, err := s.bot.session.ChannelFileSend(s.channel, name, bytes.NewReader(file))
	endTime := time.Now()
	if err != nil {
		log.Printf("%s failed to upload file %s to Discord: %v\n", s.bot.botID, name, err)
		return "", err
	}
	log.Printf("%s uploaded file %s to Discord with ID %s in %v\n", s.bot.botID, name, data.ID, endTime.Sub(startTime))
	return data.ID, nil
}

var _ service.FileService = &DiscordFileService{}
