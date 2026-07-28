package cmd

import (
	"github.com/qa-guru/cm/selenoid"
	"github.com/spf13/cobra"
)

var selenoidStartUICmd = &cobra.Command{
	Use:   "start",
	Short: "Start Selenoid UI",
	Run: func(cmd *cobra.Command, args []string) {
		startImpl(uiConfigDir, uiPort, func(lc *selenoid.Lifecycle) error {
			return lc.StartUI()
		}, force)
	},
}
