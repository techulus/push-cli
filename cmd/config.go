package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/techulus/push-cli/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

var setKeyCmd = &cobra.Command{
	Use:   "set-key <api-key>",
	Short: "Save your Push API key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := strings.TrimSpace(args[0])
		if key == "" {
			fmt.Fprintln(os.Stderr, "API key cannot be empty")
			os.Exit(1)
		}
		if err := config.SetAPIKey(key); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving API key: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("API key saved successfully")
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		key := config.GetAPIKey()
		if key == "" {
			fmt.Println("No API key configured. Run: push config set-key <api-key>")
			return
		}
		fmt.Printf("API Key: %s\n", config.MaskedAPIKey())
	},
}

func init() {
	configCmd.AddCommand(setKeyCmd)
	configCmd.AddCommand(showCmd)
	rootCmd.AddCommand(configCmd)
}
