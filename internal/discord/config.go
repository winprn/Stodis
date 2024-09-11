package discord

type BotConfig struct {
	Token string
}

func NewBotConfig(token string) *BotConfig {
	return &BotConfig{
		Token: token,
	}
}
