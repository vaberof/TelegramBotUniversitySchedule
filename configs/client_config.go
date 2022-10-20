package configs

import "time"

type HttpClientConfig struct {
	Host    string
	Timeout time.Duration
}

func NewHttpClientConfig(host string, timeout time.Duration) *HttpClientConfig {
	return &HttpClientConfig{
		Host:    host,
		Timeout: timeout,
	}
}
