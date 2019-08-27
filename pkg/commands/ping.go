// Package commands contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits
package commands

import (
	"strconv"

	"github.com/takama/grpc/pkg/boot"
	"github.com/takama/grpc/pkg/client"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping client/server connection",
	Long:  `This command provide ping service via gRPC.`,
	// nolint: unparam
	Run: func(cmd *cobra.Command, args []string) {
		cfg, log := boot.Setup()
		// nolint: errcheck
		defer log.Sync()

		// Read counts, initial value is 1
		count := 1
		countParam := cmd.Flag("count").Value.String()
		if v, err := strconv.Atoi(countParam); err == nil {
			count = v
		}

		// Runs the domain checker
		if err := client.Ping(
			&cfg.Client, log,
			cmd.Flag("message").Value.String(), count,
			boot.PrepareDialOptions(
				cfg.Client.Host, cfg.Client.Insecure,
				cfg.Client.WaitForReady, cfg.Client.BackoffDelay,
			)...,
		); err != nil {
			log.Fatal("Get info error", zap.Error(err))
		}
	},
}

func init() {
	clientCmd.AddCommand(pingCmd)

	pingCmd.PersistentFlags().String("message", "Hello", "Specify message to ping")
	pingCmd.PersistentFlags().Int("count", 1, "Count of messages to ping")
}
