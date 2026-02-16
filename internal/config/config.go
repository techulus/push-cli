package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var osExit = os.Exit

const (
	configDir  = ".push"
	configFile = "config"
	configType = "yaml"
)

func Init() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
		osExit(1)
	}

	dir := filepath.Join(home, configDir)
	viper.AddConfigPath(dir)
	viper.SetConfigName(configFile)
	viper.SetConfigType(configType)
	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
			osExit(1)
		}
	}
}

func SetAPIKey(key string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("finding home directory: %w", err)
	}

	dir := filepath.Join(home, configDir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	viper.Set("api_key", key)
	configPath := filepath.Join(dir, configFile+".yaml")
	if err := viper.WriteConfigAs(configPath); err != nil {
		return err
	}
	return os.Chmod(configPath, 0600)
}

func GetAPIKey() string {
	return viper.GetString("api_key")
}

func MaskedAPIKey() string {
	key := GetAPIKey()
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
