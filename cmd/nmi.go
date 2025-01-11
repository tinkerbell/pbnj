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
	sendNMICmd = &cobra.Command{
		Use:   "nmi",
		Short: "Send NMI",
		Long:  `Tell the target BMC to send an NMI to the system`,
		Run: func(_ *cobra.Command, _ []string) {
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

			err = v1Client.SendNMI(ctx, client, &v1.SendNMIRequest{
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
			})
			if err != nil {
				logger.Error(err, "error calling")
				os.Exit(1)
			}
		},
	}
)

func init() {
	diagnosticCmd.AddCommand(sendNMICmd)
}
