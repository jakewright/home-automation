package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jakewright/home-automation/tools/bolt/pkg/service"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

var (
	stopCmd = &cobra.Command{
		Use:   "stop [service.foo] [service.bar]...",
		Short: "stop a service",
		Run: func(cmd *cobra.Command, args []string) {
			all, err := cmd.Flags().GetBool("all")
			if err != nil {
				output.Fatal("Failed to parse all flag: %v", err)
			}

			if all {
				if err := service.StopAll(); err != nil {
					output.Fatal("Failed to stop services: %v", err)
				}

				return
			}

			if err := service.Stop(args); err != nil {
				output.Fatal("Failed to stop services: %v", err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().Bool("all", false, "stop all services")
}
