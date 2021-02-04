package cmd

import (
	"github.com/spf13/cobra"
)

// clientCmd represents the server command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Run PBnJ client",
	Long:  `Run PBnJ client for server interactions.`,
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
