package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	botID string
	session *discordgo.Session
}

func NewBot(botID, token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		botID: botID,
		session: session,
	}, nil
}
