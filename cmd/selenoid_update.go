package cmd

import (
	"github.com/qa-guru/cm/selenoid"
	"github.com/spf13/cobra"
)

var selenoidUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Selenoid (download latest Selenoid, configure and start)",
	Run: func(cmd *cobra.Command, args []string) {
		startImpl(configDir, port, func(lc *selenoid.Lifecycle) error {
			return lc.Start()
		}, true)
	},
}
