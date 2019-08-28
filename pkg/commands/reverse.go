// Package commands contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits, dupl
package commands

import (
	"strconv"

	"github.com/takama/grpc/pkg/boot"
	"github.com/takama/grpc/pkg/client"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// reverseCmd represents a reverse ping command
var reverseCmd = &cobra.Command{
	Use:   "reverse",
	Short: "Reverse ping client/serverice:services connection",
	Long:  `This command provide reverse ping between gRPC services.`,
	// nolint: unparam
	Run: func(cmd *cobra.Command, args []string) {
		cfg, log := boot.Setup()
		// nolint: errcheck
		defer log.Sync()

		// Ping counts, initial value is 1
		count := 1
		countParam := cmd.Flag("count").Value.String()
		if v, err := strconv.Atoi(countParam); err == nil {
			count = v
		}

		// Create new client
		cl, err := client.New(
			&cfg.Client, log,
			boot.PrepareDialOptions(
				cfg.Client.Host, cfg.Client.Insecure,
				cfg.Client.WaitForReady, cfg.Client.BackOffDelay,
			)...,
		)
		if err != nil {
			log.Fatal("Get connection error", zap.Error(err))
		}
		// Reverse ping command
		if err := cl.Reverse(cmd.Flag("message").Value.String(), count); err != nil {
			log.Fatal("Reverse ping error", zap.Error(err))
		}
		cl.Shutdown()
	},
}

func init() {
	clientCmd.AddCommand(reverseCmd)

	reverseCmd.PersistentFlags().String("message", "Hello", "Specify message for reverse ping")
	reverseCmd.PersistentFlags().Int("count", 1, "Count of messages for reverse ping")
}
