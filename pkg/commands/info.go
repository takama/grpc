// Package commands contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits
package commands

import (
	"context"
	"time"

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

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		// Create new client
		cl, err := client.New(ctx, &cfg.Client, log)
		if err != nil {
			log.Fatal("Get connection error", zap.Error(err))
		}
		// Runs the domain checker
		if err := cl.Info(ctx); err != nil {
			log.Fatal("Get info error", zap.Error(err))
		}
		if err := cl.Shutdown(ctx); err != nil {
			log.Fatal("Close connection error", zap.Error(err))
		}
	},
}

func init() {
	clientCmd.AddCommand(infoCmd)
}
