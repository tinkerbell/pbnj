package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	v1Client "github.com/tinkerbell/pbnj/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	bmcaddress string
	bmcpass    string
	bmcuser    string
	bmcvendor  string

	screenshotCmd = &cobra.Command{
		Use:   "screenshot",
		Short: "Take a screenshot",
		Long:  `Take a screenshot of a target BMC`,
		Run: func(cmd *cobra.Command, args []string) {
			var opts []grpc.DialOption
			ctx := context.Background()

			logger := defaultLogger(logLevel)

			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
			conn, err := grpc.Dial("localhost:"+port, opts...)
			if err != nil {
				logger.Error(err, "fail to dial server")
				os.Exit(1)
			}
			defer conn.Close()
			client := v1.NewDiagnosticClient(conn)

			resp, err := v1Client.Screenshot(ctx, client, &v1.ScreenshotRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_DirectAuthn{
						DirectAuthn: &v1.DirectAuthn{
							Host: &v1.Host{
								Host: bmcaddress,
							},
							Username: bmcuser,
							Password: bmcpass,
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: bmcvendor,
				},
			})
			if err != nil {
				logger.Error(err, "error calling")
				os.Exit(1)
			}

			logger.Info("resp", "resp", []interface{}{resp})
		},
	}
)

func init() {
	diagnosticCmd.PersistentFlags().StringVar(&bmcaddress, "bmcaddress", "", "bmc address")
	diagnosticCmd.PersistentFlags().StringVar(&bmcuser, "bmcuser", "", "bmc user")
	diagnosticCmd.PersistentFlags().StringVar(&bmcpass, "bmcpass", "", "bmc password")
	diagnosticCmd.PersistentFlags().StringVar(&bmcvendor, "bmcvendor", "", "bmc vendor")

	diagnosticCmd.AddCommand(screenshotCmd)
}
