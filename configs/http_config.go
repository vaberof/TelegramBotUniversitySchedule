package configs

import "time"

type HttpClientConfig struct {
	Timeout time.Duration
}

func NewHttpClientConfig(timeout time.Duration) *HttpClientConfig {
	return &HttpClientConfig{Timeout: timeout}
}
