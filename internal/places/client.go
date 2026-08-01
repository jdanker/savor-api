package places

import (
	"net/http"
	"time"
)

type Client struct {
	apikey  string
	baseUrl string
	http    *http.Client
}

func New(apikey string) *Client {
	return &Client{
		apikey:  apikey,
		baseUrl: "https://places.googleapis.com/v1",
		http:    &http.Client{Timeout: 5 * time.Second},
	}
}
