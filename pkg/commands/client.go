// Package commands contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits
package commands

import (
	"github.com/takama/grpc/pkg/config"
	"github.com/takama/grpc/pkg/helper"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client commands",
	Long:  `These commands provides ability to check connections using gRPC client`,
}

func init() {
	RootCmd.AddCommand(clientCmd)

	clientCmd.PersistentFlags().String("host", config.ClientServiceName, "Client service host")
	clientCmd.PersistentFlags().Int("port", config.DefaultClientPort, "Client service port")
	helper.LogF("Flag error", viper.BindPFlag("client.host", clientCmd.PersistentFlags().Lookup("host")))
	helper.LogF("Flag error", viper.BindPFlag("client.port", clientCmd.PersistentFlags().Lookup("port")))
	helper.LogF("Env error", viper.BindEnv("client.host"))
	helper.LogF("Env error", viper.BindEnv("client.port"))
}
