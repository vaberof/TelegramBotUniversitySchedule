package integration

import (
	"github.com/go-resty/resty/v2"
	"github.com/vaberof/TelegramBotUniversitySchedule/configs"
)

type HttpClient struct {
	client *resty.Client
}

func NewHttpClient(cfg *configs.HttpClientConfig) *HttpClient {
	return &HttpClient{
		client: resty.New().SetTimeout(cfg.Timeout),
	}
}
