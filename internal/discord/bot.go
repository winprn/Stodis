package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
}

func NewBot(config *BotConfig) (*Bot, error) {
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Session: session,
	}, nil
}
