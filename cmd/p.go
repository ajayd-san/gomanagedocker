/*
Copyright © 2024 Saivenkat Ajay D. <ajayds2001@gmail.com>
*/
package cmd

import (
	"github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/ajayd-san/gomanagedocker/tui"
	"github.com/spf13/cobra"
)

// pCmd represents the p command
var pCmd = &cobra.Command{
	Use:   "p",
	Short: "Manage Podman objects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.StartTUI(debug, types.Podman)
	},
}

func init() {
	rootCmd.AddCommand(pCmd)
}
