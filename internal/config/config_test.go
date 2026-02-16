package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestMaskedAPIKey_Long(t *testing.T) {
	viper.Set("api_key", "abcd1234efgh5678")
	got := MaskedAPIKey()
	want := "abcd...5678"
	if got != want {
		t.Errorf("MaskedAPIKey() = %q, want %q", got, want)
	}
}

func TestMaskedAPIKey_Short(t *testing.T) {
	viper.Set("api_key", "short")
	got := MaskedAPIKey()
	if got != "****" {
		t.Errorf("MaskedAPIKey() = %q, want %q", got, "****")
	}
}

func TestMaskedAPIKey_Empty(t *testing.T) {
	viper.Set("api_key", "")
	got := MaskedAPIKey()
	if got != "****" {
		t.Errorf("MaskedAPIKey() = %q, want %q", got, "****")
	}
}

func TestMaskedAPIKey_ExactlyEight(t *testing.T) {
	viper.Set("api_key", "12345678")
	got := MaskedAPIKey()
	if got != "****" {
		t.Errorf("MaskedAPIKey() = %q, want %q", got, "****")
	}
}

func TestMaskedAPIKey_NineChars(t *testing.T) {
	viper.Set("api_key", "123456789")
	got := MaskedAPIKey()
	want := "1234...6789"
	if got != want {
		t.Errorf("MaskedAPIKey() = %q, want %q", got, want)
	}
}

func TestInit_MalformedConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cfgBase, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("UserConfigDir() error: %v", err)
	}
	cfgDir := filepath.Join(cfgBase, "push")
	os.MkdirAll(cfgDir, 0700)
	os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte(":::invalid yaml"), 0600)

	viper.Reset()

	exitCalled := false
	origExit := osExit
	osExit = func(code int) {
		exitCalled = true
	}
	defer func() { osExit = origExit }()

	Init()

	if !exitCalled {
		t.Error("expected os.Exit to be called for malformed config")
	}
}

func TestSetAPIKey_FilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	viper.Reset()

	if err := SetAPIKey("test-key"); err != nil {
		t.Fatalf("SetAPIKey() error: %v", err)
	}

	cfgBase, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("UserConfigDir() error: %v", err)
	}
	configPath := filepath.Join(cfgBase, "push", "config.yaml")
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Stat(%q) error: %v", configPath, err)
	}

	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("config file permissions = %04o, want 0600", perm)
	}
}

func TestGetAPIKey(t *testing.T) {
	viper.Set("api_key", "my-test-key")
	got := GetAPIKey()
	if got != "my-test-key" {
		t.Errorf("GetAPIKey() = %q, want %q", got, "my-test-key")
	}
}
