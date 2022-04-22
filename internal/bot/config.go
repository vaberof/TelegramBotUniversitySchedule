package bot

type Config struct {
	Token string `yaml:"token"`
}

// NewConfig returns pointer to Config.
func NewConfig() *Config {
	return &Config{}
}