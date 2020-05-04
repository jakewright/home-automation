package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jakewright/home-automation/tools/bolt/pkg/service"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

var (
	buildCmd = &cobra.Command{
		Use:   "build [service.foo] [service.bar]...",
		Short: "build a service",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := service.Build(args); err != nil {
				output.Fatal("Failed to build: %v", err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(buildCmd)
}
