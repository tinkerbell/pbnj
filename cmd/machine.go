package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/packethost/pkg/log/logr"
	"github.com/spf13/cobra"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	v1Client "github.com/tinkerbell/pbnj/client"
	"google.golang.org/grpc"
)

// machineCmd represents the server command
var machineCmd = &cobra.Command{
	Use:   "machine",
	Short: "Run PBnJ client machine actions",
	Long:  `Run PBnJ client for machine actions.`,
	Run: func(cmd *cobra.Command, args []string) {
		var opts []grpc.DialOption
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		logger, zlog, err := logr.NewPacketLogr(
			logr.WithServiceName("github.com/tinkerbell/pbnj"),
			logr.WithLogLevel(logLevel),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		defer zlog.Sync() // nolint

		opts = append(opts, grpc.WithInsecure())
		conn, err := grpc.Dial("localhost:"+port, opts...)
		if err != nil {
			logger.V(0).Error(err, "fail to dial server")
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
			logger.V(0).Error(err, "error calling")
			os.Exit(1)
		}

		logger.V(0).Info("resp", "resp", []interface{}{resp})

	},
}

func init() {
	machineCmd.PersistentFlags().StringVar(&port, "port", "50051", "server port (default is 50051")
	clientCmd.AddCommand(machineCmd)
}
