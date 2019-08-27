// Package commands contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits
package commands

import (
	"github.com/takama/grpc/pkg/boot"
	"github.com/takama/grpc/pkg/client"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information about service",
	Long: `This command provide service information via gRPC.
Use --file or -f to specify JSON data file with batch requests`,
	// nolint: unparam
	Run: func(cmd *cobra.Command, args []string) {
		cfg, log := boot.Setup()
		// nolint: errcheck
		defer log.Sync()

		// Runs the domain checker
		if err := client.Info(
			&cfg.Client, log,
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
	clientCmd.AddCommand(infoCmd)
}
