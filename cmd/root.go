/*
Copyright © 2024 Saivenkat Ajay D. <ajayds2001@gmail.com>
*/
package cmd

import (
	_ "embed"
	"errors"
	"os"
	"strings"

	"github.com/ajayd-san/gomanagedocker/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//go:embed defaultConfig.yaml
	defaultConfig string
	debug         bool
	rootCmd       = &cobra.Command{
		Use:     "gmd",
		Short:   "TUI to manage docker objects",
		Long:    `The Definitive TUI to manage docker objects with ease.`,
		Version: "1.1.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tui.StartTUI(debug)
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	readConfig()

	rootCmd.Flags().BoolVar(&debug, "debug", false, "Send logs to ./gmd_debug.log")
}

func readConfig() {
	//read config file
	configPath, err := os.UserConfigDir()

	if err != nil {
		configPath, err = os.UserHomeDir()

		if err != nil {
			cobra.CheckErr(err)
		}
	}

	viper.AddConfigPath(configPath + "/gomanagedocker")
	viper.SetConfigName("gomanagedocker")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, new(viper.ConfigFileNotFoundError)) {
			err := viper.ReadConfig(strings.NewReader(defaultConfig))
			if err != nil {
				panic(err)
			}
		} else {
			cobra.CheckErr(err)
		}
	}
}
