package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tinkerbell/pbnj/cmd/zaplog"
)

// clientCmd represents the server command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Run PBnJ client",
	Long:  `Run PBnJ client for server interactions.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		logger, zlog, err := zaplog.RegisterLogger(
			zaplog.WithServiceName("github.com/tinkerbell/pbnj"),
			zaplog.WithLogLevel(logLevel),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		defer zlog.Sync() // nolint

		logger.V(0).Info("not implemented")

	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
