// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits
package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/takama/grpc/pkg/config"
	"github.com/takama/grpc/pkg/helper"
	"github.com/takama/grpc/pkg/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Service short description",
	Long:  `Service long description`,
}

// Run adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Run() {
	helper.LogF("Service bootstrap error", RootCmd.Execute())
}

func init() {
	viper.SetEnvPrefix(config.ServiceName)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	viper.SetConfigType("json")
	viper.SetConfigFile(config.DefaultConfigPath)
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config/default.conf)")
	RootCmd.PersistentFlags().Int("log-level", int(config.DefaultLoggerLevel), "Logger level (0 - debug, 1 - info, ...)")
	RootCmd.PersistentFlags().String("log-format", logger.TextFormatter.String(), "Logger format: txt, json")
	helper.LogF("Flag error",
		viper.BindPFlag("logger.level", RootCmd.PersistentFlags().Lookup("log-level")))
	helper.LogF("Flag error",
		viper.BindPFlag("logger.format", RootCmd.PersistentFlags().Lookup("log-format")))
	helper.LogF("Env error", viper.BindEnv("logger.level"))
	helper.LogF("Env error", viper.BindEnv("logger.format"))
}

// initConfig reads in config file
func initConfig() {
	// enable ability to specify config file via flag or via env
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else
	// Check for env variable with config path
	if cfgPath := os.Getenv(
		strings.ToUpper(strings.Replace(config.ServiceName, "-", "_", -1)) + "_CONFIG_PATH",
	); cfgPath != "" {
		viper.SetConfigFile(cfgPath)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
