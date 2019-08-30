// Package commands contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits, dupl
package commands

import (
	"context"
	"strconv"
	"time"

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

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
		defer cancel()

		// Create new client
		cl, err := client.New(ctx, &cfg.Client, log)
		if err != nil {
			log.Fatal("Get connection error", zap.Error(err))
		}
		// Reverse ping command
		if err := cl.Reverse(cl.Content(), cmd.Flag("message").Value.String(), count); err != nil {
			log.Fatal("Reverse ping error", zap.Error(err))
		}
		if err := cl.Shutdown(ctx); err != nil {
			log.Fatal("Close connection error", zap.Error(err))
		}
	},
}

func init() {
	clientCmd.AddCommand(reverseCmd)

	reverseCmd.PersistentFlags().String("message", "Hello", "Specify message for reverse ping")
	reverseCmd.PersistentFlags().Int("count", 1, "Count of messages for reverse ping")
}
