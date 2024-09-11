package discord

import (
	"os"
	"testing"
)

func Test_DiscordServiceUploadFile(t *testing.T) {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	bot, _ := NewBot(NewBotConfig(token))
	service := NewDiscordFileService(bot, "1278013883973632071")

	service.UploadFile([]byte("test"), "test")
}
