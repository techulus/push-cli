package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/techulus/push-cli/internal/api"
	"github.com/techulus/push-cli/internal/config"
)

var validSounds = []string{
	"default", "arcade", "correct", "fail", "harp", "reveal",
	"bubble", "doorbell", "flute", "money", "scifi", "clear",
	"elevator", "guitar", "pop",
}

func readBodyFromStdinOrFlag(cmd *cobra.Command) (string, error) {
	body, _ := cmd.Flags().GetString("body")

	if body == "-" || body == "" {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return "", fmt.Errorf("reading stdin: %w", err)
			}
			trimmed := strings.TrimSpace(string(data))
			if trimmed == "" {
				return "", fmt.Errorf("body is required (use --body flag or pipe via stdin)")
			}
			return trimmed, nil
		}
	}

	if body == "" || body == "-" {
		return "", fmt.Errorf("body is required (use --body flag or pipe via stdin)")
	}

	return body, nil
}

func buildNotifyRequest(cmd *cobra.Command) (api.NotifyRequest, error) {
	title, _ := cmd.Flags().GetString("title")
	body, err := readBodyFromStdinOrFlag(cmd)
	if err != nil {
		return api.NotifyRequest{}, err
	}

	sound, _ := cmd.Flags().GetString("sound")
	if sound != "" {
		valid := false
		for _, s := range validSounds {
			if s == sound {
				valid = true
				break
			}
		}
		if !valid {
			return api.NotifyRequest{}, fmt.Errorf("invalid sound %q, valid sounds: %s", sound, strings.Join(validSounds, ", "))
		}
	}

	channel, _ := cmd.Flags().GetString("channel")
	link, _ := cmd.Flags().GetString("link")
	image, _ := cmd.Flags().GetString("image")
	timeSensitive, _ := cmd.Flags().GetBool("time-sensitive")

	return api.NotifyRequest{
		Title:         title,
		Body:          body,
		Sound:         sound,
		Channel:       channel,
		Link:          link,
		Image:         image,
		TimeSensitive: timeSensitive,
	}, nil
}

func newAPIClient() *api.Client {
	key := config.GetAPIKey()
	if key == "" {
		fmt.Fprintln(os.Stderr, "No API key configured. Run: push config set-key <api-key>")
		os.Exit(1)
	}
	return api.NewClient(key)
}

func addNotifyFlags(cmd *cobra.Command) {
	cmd.Flags().String("title", "", "Notification title (required)")
	cmd.Flags().String("body", "", "Notification body (use '-' to read from stdin)")
	cmd.Flags().String("sound", "", "Notification sound")
	cmd.Flags().String("channel", "", "Notification channel")
	cmd.Flags().String("link", "", "URL to open when notification is tapped")
	cmd.Flags().String("image", "", "Image URL for the notification")
	cmd.Flags().Bool("time-sensitive", false, "Mark as time-sensitive")
	cmd.MarkFlagRequired("title")
}

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Send a push notification",
	Run: func(cmd *cobra.Command, args []string) {
		req, err := buildNotifyRequest(cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		client := newAPIClient()
		resp, err := client.Notify(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(resp)
	},
}

func init() {
	addNotifyFlags(notifyCmd)
	rootCmd.AddCommand(notifyCmd)
}
