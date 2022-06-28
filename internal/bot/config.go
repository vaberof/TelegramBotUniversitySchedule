package bot

type Config struct {
	Token string `yaml:"token"`
}

func NewConfig() *Config {
	return &Config{}
}
