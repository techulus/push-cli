package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var baseURL = "https://push.techulus.com/api/v1"

type NotifyRequest struct {
	Title         string `json:"title"`
	Body          string `json:"body"`
	Sound         string `json:"sound,omitempty"`
	Channel       string `json:"channel,omitempty"`
	Link          string `json:"link,omitempty"`
	Image         string `json:"image,omitempty"`
	TimeSensitive bool   `json:"timeSensitive,omitempty"`
}

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) Notify(req NotifyRequest) (string, error) {
	return c.post("/notify", req)
}

func (c *Client) NotifyAsync(req NotifyRequest) (string, error) {
	return c.post("/notify-async", req)
}

func (c *Client) NotifyGroup(groupID string, req NotifyRequest) (string, error) {
	return c.post("/notify/group/"+url.PathEscape(groupID), req)
}

func (c *Client) post(path string, payload interface{}) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", baseURL+path, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}
