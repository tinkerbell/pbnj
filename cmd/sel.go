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
	clearSystemEventLogcommand = &cobra.Command{
		Use:   "selclear",
		Short: "Clear the System Event Log",
		Long:  `Clear the System Event Log of the target BMC`,
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

			taskClient := v1.NewTaskClient(conn)

			resp, err := v1Client.ClearSystemEventLog(ctx, client, taskClient, &v1.ClearSystemEventLogRequest{
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

	getSystemEventLogcommand = &cobra.Command{
		Use:   "selget",
		Short: "Get the System Event Log",
		Long: `Get the System Event Log of the target BMC.
		Includes the following information:
		ID, Timestamp, Description, Message`,
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

			resp, err := v1Client.SystemEventLog(ctx, client, &v1.SystemEventLogRequest{
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

	getSystemEventLogRawcommand = &cobra.Command{
		Use:   "selgetraw",
		Short: "Get the Raw System Event Log",
		Long:  `Get the Raw System Event Log of the target BMC`,
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

			resp, err := v1Client.SystemEventLogRaw(ctx, client, &v1.SystemEventLogRawRequest{
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
	diagnosticCmd.AddCommand(clearSystemEventLogcommand)
	diagnosticCmd.AddCommand(getSystemEventLogcommand)
	diagnosticCmd.AddCommand(getSystemEventLogRawcommand)
}
