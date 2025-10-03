package service

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	resty *resty.Client
}

func NewClient(baseURL string, timeout time.Duration, defaultHeaders map[string]string) *Client {
	c := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(timeout)

	// 设置默认 Header
	for k, v := range defaultHeaders {
		c.SetHeader(k, v)
	}

	return &Client{
		resty: c,
	}
}
