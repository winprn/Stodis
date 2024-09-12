package discord

import (
	"os"
	"testing"
)

func Test_DiscordServiceUploadFile(t *testing.T) {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	bot, _ := NewBot(NewBotConfig(token))
	service := NewDiscordHandler(bot, "1278013883973632071")

	if _, err := service.UploadFile([]byte("test"), "test"); err != nil {
		t.Error(err)
	}
}
