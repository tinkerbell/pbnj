package cmd

import (
	"github.com/spf13/cobra"
)

var (
	bmcaddress string
	bmcpass    string
	bmcuser    string
	bmcvendor  string
)

// diagnosticCmd represents the diagnostic command.
var diagnosticCmd = &cobra.Command{
	Use:   "diagnostic",
	Short: "Run PBnJ client",
	Long:  `Run PBnJ client for server interactions.`,
}

func init() {
	diagnosticCmd.PersistentFlags().StringVar(&bmcaddress, "bmcaddress", "", "bmc address")
	diagnosticCmd.PersistentFlags().StringVar(&bmcuser, "bmcuser", "", "bmc user")
	diagnosticCmd.PersistentFlags().StringVar(&bmcpass, "bmcpass", "", "bmc password")
	diagnosticCmd.PersistentFlags().StringVar(&bmcvendor, "bmcvendor", "", "bmc vendor")

	clientCmd.AddCommand(diagnosticCmd)
}
