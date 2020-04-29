// Package commands contains global variables with configs and commands.
// nolint: gochecknoglobals, gochecknoinits
package commands

import (
	"github.com/takama/grpc/pkg/config"
	"github.com/takama/grpc/pkg/helper"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clientCmd represents the client command.
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client commands",
	Long:  `These commands provides ability to check connections using gRPC client`,
}

func init() {
	RootCmd.AddCommand(clientCmd)

	clientCmd.PersistentFlags().String("scheme", config.ServiceName, "Client target scheme")
	clientCmd.PersistentFlags().String("host", config.ClientServiceName, "Client service host")
	clientCmd.PersistentFlags().String("balancer", config.DefaultClientBalancer, "Client service balancer")
	helper.LogF("Flag error", viper.BindPFlag("client.scheme", clientCmd.PersistentFlags().Lookup("scheme")))
	helper.LogF("Flag error", viper.BindPFlag("client.host", clientCmd.PersistentFlags().Lookup("host")))
	helper.LogF("Flag error", viper.BindPFlag("client.balancer", clientCmd.PersistentFlags().Lookup("balancer")))
	helper.LogF("Env error", viper.BindEnv("client.scheme"))
	helper.LogF("Env error", viper.BindEnv("client.host"))
	helper.LogF("Env error", viper.BindEnv("client.balancer"))
}
