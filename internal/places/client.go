package places

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
		// Own timeout below the server's 10s WriteTimeout so a slow Google
		// call fails inside our request budget instead of at it.
		http: &http.Client{Timeout: 5 * time.Second},
	}
}

// localizedText is Google's {text, languageCode} wrapper, shared by every endpoint.
type localizedText struct {
	Text string `json:"text"`
}

// do is the shared authenticated-request helper: build request → attach key +
// field mask → execute → status check → decode into dest.
//
// Auth and the field mask go in headers (X-Goog-*), not the URL — keeps the key
// out of logs and works uniformly for GET and POST.
func (c *Client) do(ctx context.Context, method, path, fieldMask string, body, dest any) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseUrl+path, bodyReader)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", c.apikey)
	if fieldMask != "" {
		req.Header.Set("X-Goog-FieldMask", fieldMask)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("call google places: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read a bounded slice of the error body — enough to debug,
		// impossible for a hostile/huge response to balloon memory.
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<10))
		return fmt.Errorf("google places %s %s: status %d: %s", method, path, resp.StatusCode, msg)
	}

	if dest == nil {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
