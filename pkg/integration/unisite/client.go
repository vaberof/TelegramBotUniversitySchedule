package integration

import (
	"github.com/go-resty/resty/v2"
	"github.com/vaberof/TelegramBotUniversitySchedule/configs"
)

type HttpClient struct {
	client *resty.Client
	host   string
}

func NewHttpClient(config *configs.HttpClientConfig) *HttpClient {
	return &HttpClient{
		host:   config.Host,
		client: resty.New().SetTimeout(config.Timeout),
	}
}
