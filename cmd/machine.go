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

// machineCmd represents the server command.
var machineCmd = &cobra.Command{
	Use:   "machine",
	Short: "Run PBnJ client machine actions",
	Long:  `Run PBnJ client for machine actions.`,
	Run: func(_ *cobra.Command, _ []string) {
		var opts []grpc.DialOption
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		logger := defaultLogger(logLevel)

		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		conn, err := grpc.NewClient("localhost:"+port, opts...)
		if err != nil {
			logger.Error(err, "fail to dial server")
			os.Exit(1)
		}
		defer conn.Close()
		client := v1.NewMachineClient(conn)
		taskClient := v1.NewTaskClient(conn)

		resp, err := v1Client.MachinePower(ctx, client, taskClient, &v1.PowerRequest{
			Authn: &v1.Authn{
				Authn: &v1.Authn_DirectAuthn{
					DirectAuthn: &v1.DirectAuthn{
						Host: &v1.Host{
							Host: "127.0.0.1",
						},
						Username: "ADMIN",
						Password: "ADMIN",
					},
				},
			},
			Vendor: &v1.Vendor{
				Name: "HP",
			},
			PowerAction: v1.PowerAction_POWER_ACTION_STATUS,
		})
		if err != nil {
			logger.Error(err, "error calling")
			os.Exit(1)
		}

		logger.Info("resp", "resp", []interface{}{resp})
	},
}

func init() {
	machineCmd.PersistentFlags().StringVar(&port, "port", "50051", "server port (default is 50051")
	clientCmd.AddCommand(machineCmd)
}
