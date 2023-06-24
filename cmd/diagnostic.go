package cmd

import (
	"github.com/spf13/cobra"
)

// diagnosticCmd represents the diagnostic command.
var diagnosticCmd = &cobra.Command{
	Use:   "diagnostic",
	Short: "Run PBnJ client",
	Long:  `Run PBnJ client for server interactions.`,
}

func init() {
	clientCmd.AddCommand(diagnosticCmd)
}
