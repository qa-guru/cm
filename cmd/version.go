package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	gitRevision string = "HEAD"
	buildStamp  string = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "Git Revision: %s\n", gitRevision)
		fmt.Fprintf(cmd.OutOrStdout(), "UTC Build Time: %s\n", buildStamp)
	},
}
