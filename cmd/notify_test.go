package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	addNotifyFlags(cmd)
	return cmd
}

func TestBuildNotifyRequest_ValidSound(t *testing.T) {
	cmd := newTestCmd()
	cmd.SetArgs([]string{"--title", "Test", "--body", "Hello", "--sound", "arcade"})
	cmd.Execute()

	req, err := buildNotifyRequest(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Sound != "arcade" {
		t.Errorf("expected sound arcade, got %s", req.Sound)
	}
}

func TestBuildNotifyRequest_InvalidSound(t *testing.T) {
	cmd := newTestCmd()
	cmd.SetArgs([]string{"--title", "Test", "--body", "Hello", "--sound", "invalid"})
	cmd.Execute()

	_, err := buildNotifyRequest(cmd)
	if err == nil {
		t.Fatal("expected error for invalid sound")
	}
}

func TestBuildNotifyRequest_AllFlags(t *testing.T) {
	cmd := newTestCmd()
	cmd.SetArgs([]string{
		"--title", "Test",
		"--body", "Hello",
		"--sound", "pop",
		"--channel", "alerts",
		"--link", "https://example.com",
		"--image", "https://example.com/img.png",
		"--time-sensitive",
	})
	cmd.Execute()

	req, err := buildNotifyRequest(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Title != "Test" {
		t.Errorf("expected title Test, got %s", req.Title)
	}
	if req.Body != "Hello" {
		t.Errorf("expected body Hello, got %s", req.Body)
	}
	if req.Sound != "pop" {
		t.Errorf("expected sound pop, got %s", req.Sound)
	}
	if req.Channel != "alerts" {
		t.Errorf("expected channel alerts, got %s", req.Channel)
	}
	if req.Link != "https://example.com" {
		t.Errorf("expected link, got %s", req.Link)
	}
	if req.Image != "https://example.com/img.png" {
		t.Errorf("expected image, got %s", req.Image)
	}
	if !req.TimeSensitive {
		t.Error("expected time-sensitive to be true")
	}
}

func TestBuildNotifyRequest_NoBody(t *testing.T) {
	cmd := newTestCmd()
	cmd.SetArgs([]string{"--title", "Test"})
	cmd.Execute()

	_, err := buildNotifyRequest(cmd)
	if err == nil {
		t.Fatal("expected error when body is missing")
	}
}

func TestReadBodyFromStdin_Empty(t *testing.T) {
	r, w, _ := os.Pipe()
	w.WriteString("   \n  ")
	w.Close()

	origStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	cmd := newTestCmd()
	cmd.SetArgs([]string{"--title", "Test"})
	cmd.Execute()

	_, err := readBodyFromStdinOrFlag(cmd)
	if err == nil {
		t.Fatal("expected error for empty stdin body")
	}
	if !strings.Contains(err.Error(), "body is required") {
		t.Errorf("expected 'body is required' error, got: %v", err)
	}
}

func TestValidSounds(t *testing.T) {
	expected := []string{
		"default", "arcade", "correct", "fail", "harp", "reveal",
		"bubble", "doorbell", "flute", "money", "scifi", "clear",
		"elevator", "guitar", "pop",
	}
	if len(validSounds) != len(expected) {
		t.Fatalf("expected %d sounds, got %d", len(expected), len(validSounds))
	}
	for i, s := range expected {
		if validSounds[i] != s {
			t.Errorf("expected sound %q at index %d, got %q", s, i, validSounds[i])
		}
	}
}
