// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits, unparam
package commands

import (
	"github.com/takama/grpc/pkg/config"
	"github.com/takama/grpc/pkg/helper"
	"github.com/takama/grpc/pkg/service"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Listen and handle requests including health/ready checks",
	Long: `This command prepare the service for handling
of the requests to the service.
Also there are setup a health check and a readiness check
which should observe a liveness/readiness of registered modules`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.New()
		helper.LogF("Loading config error", err)
		// Runs the service
		helper.LogF("Service start error", service.Run(cfg))
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().Int("server-port", config.DefaultServerPort, "Service listening port number")
	serveCmd.PersistentFlags().Int("info-port", config.DefaultInfoPort, "Health port number")
	serveCmd.PersistentFlags().Bool("info-statistics", config.DefaultInfoStatistics, "Collect statistics information")
	serveCmd.PersistentFlags().Int("grace-period", config.DefaultGracePeriod, "Service termination grace period")
	helper.LogF("Flag error",
		viper.BindPFlag("server.port", serveCmd.PersistentFlags().Lookup("server-port")))
	helper.LogF("Flag error",
		viper.BindPFlag("info.port", serveCmd.PersistentFlags().Lookup("info-port")))
	helper.LogF("Flag error",
		viper.BindPFlag("info.statistics", serveCmd.PersistentFlags().Lookup("info-statistics")))
	helper.LogF("Flag error",
		viper.BindPFlag("system.grace.period", serveCmd.PersistentFlags().Lookup("grace-period")))
	helper.LogF("Env error", viper.BindEnv("server.port"))
	helper.LogF("Env error", viper.BindEnv("info.port"))
	helper.LogF("Env error", viper.BindEnv("info.statistics"))
	helper.LogF("Env error", viper.BindEnv("system.grace.period"))
}
