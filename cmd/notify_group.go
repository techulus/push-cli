package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var notifyGroupCmd = &cobra.Command{
	Use:   "notify-group <group-id>",
	Short: "Send a push notification to a group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		groupID := args[0]

		req, err := buildNotifyRequest(cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		client := newAPIClient()
		resp, err := client.NotifyGroup(groupID, req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(resp)
	},
}

func init() {
	addNotifyFlags(notifyGroupCmd)
	rootCmd.AddCommand(notifyGroupCmd)
}
