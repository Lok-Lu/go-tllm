package hunyuan

import (
	"net/http"
)

const (
	hunyuanAPIBaseUrl              = "hunyuan.cloud.tencent.com"
	defaultEmptyMessagesLimit uint = 300
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	accessKey          string
	secretKey          string
	appid              int
	BaseURL            string
	HTTPClient         *http.Client
	EmptyMessagesLimit uint
}

func DefaultConfig(ak, sk string, appid int) (ClientConfig, error) {
	return ClientConfig{
		accessKey:  ak,
		secretKey:  sk,
		appid:      appid,
		BaseURL:    hunyuanAPIBaseUrl,
		HTTPClient: &http.Client{},

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}, nil
}

func (c ClientConfig) WithHttpClientConfig(client *http.Client) ClientConfig {
	c.HTTPClient = client
	return c
}

func (ClientConfig) String() string {
	return "<OpenAI API ClientConfig>"
}
