package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClientTimeout(t *testing.T) {
	client := NewClient("test-key")
	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("expected timeout 30s, got %v", client.httpClient.Timeout)
	}
}

func TestNotifyGroupURLEscape(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawPath != "/notify/group/my%2Fgroup%20name" {
			t.Errorf("expected escaped path /notify/group/my%%2Fgroup%%20name, got %s", r.URL.RawPath)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	origBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = origBaseURL }()

	client := NewClient("test-key")
	_, err := client.NotifyGroup("my/group name", NotifyRequest{Title: "Test", Body: "Hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNotify(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/notify" {
			t.Errorf("expected /notify, got %s", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "test-key" {
			t.Errorf("expected x-api-key test-key, got %s", r.Header.Get("x-api-key"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		var req NotifyRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Title != "Test" {
			t.Errorf("expected title Test, got %s", req.Title)
		}
		if req.Body != "Hello" {
			t.Errorf("expected body Hello, got %s", req.Body)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	origBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = origBaseURL }()

	client := NewClient("test-key")
	resp, err := client.Notify(NotifyRequest{Title: "Test", Body: "Hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != `{"success":true}` {
		t.Errorf("unexpected response: %s", resp)
	}
}

func TestNotifyAsync(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/notify-async" {
			t.Errorf("expected /notify-async, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	origBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = origBaseURL }()

	client := NewClient("test-key")
	_, err := client.NotifyAsync(NotifyRequest{Title: "Test", Body: "Hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNotifyGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/notify/group/my-group" {
			t.Errorf("expected /notify/group/my-group, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	origBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = origBaseURL }()

	client := NewClient("test-key")
	_, err := client.NotifyGroup("my-group", NotifyRequest{Title: "Test", Body: "Hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNotifyOptionalFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req NotifyRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Sound != "arcade" {
			t.Errorf("expected sound arcade, got %s", req.Sound)
		}
		if req.Channel != "alerts" {
			t.Errorf("expected channel alerts, got %s", req.Channel)
		}
		if req.Link != "https://example.com" {
			t.Errorf("expected link https://example.com, got %s", req.Link)
		}
		if req.Image != "https://example.com/img.png" {
			t.Errorf("expected image URL, got %s", req.Image)
		}
		if !req.TimeSensitive {
			t.Error("expected timeSensitive to be true")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	origBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = origBaseURL }()

	client := NewClient("test-key")
	_, err := client.Notify(NotifyRequest{
		Title:         "Test",
		Body:          "Hello",
		Sound:         "arcade",
		Channel:       "alerts",
		Link:          "https://example.com",
		Image:         "https://example.com/img.png",
		TimeSensitive: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNotifyHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"success":false,"message":"Invalid API key"}`))
	}))
	defer server.Close()

	origBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = origBaseURL }()

	client := NewClient("bad-key")
	_, err := client.Notify(NotifyRequest{Title: "Test", Body: "Hello"})
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
}

func TestOmitEmptyFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var raw map[string]interface{}
		json.NewDecoder(r.Body).Decode(&raw)

		for _, field := range []string{"sound", "channel", "link", "image", "timeSensitive"} {
			if _, exists := raw[field]; exists {
				t.Errorf("expected field %s to be omitted", field)
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	origBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = origBaseURL }()

	client := NewClient("test-key")
	_, err := client.Notify(NotifyRequest{Title: "Test", Body: "Hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
