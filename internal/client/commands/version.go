package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCommand(version, buildDate, buildCommit string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("vaultfactory CLI v%s\n", getBuildInfo(version))
			fmt.Printf("build date: %s\n", getBuildInfo(buildDate))
			fmt.Printf("build commit: %s\n", getBuildInfo(buildCommit))
		},
	}
}

func getBuildInfo(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}
