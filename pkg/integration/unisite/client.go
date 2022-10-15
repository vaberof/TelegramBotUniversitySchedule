package integration

import (
	"github.com/go-resty/resty/v2"
	"github.com/vaberof/TelegramBotUniversitySchedule/configs"
)

type HttpClient struct {
	client *resty.Client
	host   string
}

func NewHttpClient(host string, config *configs.HttpClientConfig) *HttpClient {
	return &HttpClient{
		host:   host,
		client: resty.New().SetTimeout(config.Timeout),
	}
}
