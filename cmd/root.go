/*
Copyright © 2024 Saivenkat Ajay D. <ajayds2001@gmail.com>
*/
package cmd

import (
	_ "embed"
	"os"

	"github.com/ajayd-san/gomanagedocker/tui"
	"github.com/spf13/cobra"
)

var (
	debug   bool
	rootCmd = &cobra.Command{
		Use:     "gmd",
		Short:   "TUI to manage docker objects",
		Long:    `The Definitive TUI to manage docker objects with ease.`,
		Version: "1.3.1-hotfix",
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
	rootCmd.Flags().BoolVar(&debug, "debug", false, "Send logs to ./gmd_debug.log")
}
