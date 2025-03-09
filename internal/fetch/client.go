package fetch

import (
	"net/http"
	"time"
)

var client *Client

func init() {
	client = NewClient(5 * time.Second)
}

type Client struct {
	httpClient http.Client
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}
