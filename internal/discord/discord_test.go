package discord

import (
	// "fmt"
	"os"
	"testing"
)

func Test_DiscordServiceUploadFile(t *testing.T) {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	t.Logf("token: %s", token)
	bot, _ := NewBot("bot_1", token)
	service := NewDiscordFileService(bot, "1278013883973632071")

	service.UploadFile([]byte("test"), "test")
}
