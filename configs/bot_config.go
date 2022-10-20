package configs

type BotConfig struct {
	Token string
}

func NewBotConfig(token string) *BotConfig {
	return &BotConfig{
		Token: token,
	}
}
