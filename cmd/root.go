package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/techulus/push-cli/internal/config"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "push",
	Short:   "Push by Techulus - Send push notifications from the command line",
	Version: Version,
}

func Execute() {
	config.Init()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
