package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jakewright/home-automation/tools/bolt/pkg/service"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

var (
	restartCmd = &cobra.Command{
		Use:   "restart [service.foo] [service.bar]",
		Short: "restart a service",
		Run: func(cmd *cobra.Command, args []string) {
			if err := service.Restart(args); err != nil {
				output.Fatal("Failed to restart services: %v", err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(restartCmd)
}
