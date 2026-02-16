package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var notifyAsyncCmd = &cobra.Command{
	Use:   "notify-async",
	Short: "Send a push notification asynchronously",
	Run: func(cmd *cobra.Command, args []string) {
		req, err := buildNotifyRequest(cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		client := newAPIClient()
		resp, err := client.NotifyAsync(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(resp)
	},
}

func init() {
	addNotifyFlags(notifyAsyncCmd)
	rootCmd.AddCommand(notifyAsyncCmd)
}
