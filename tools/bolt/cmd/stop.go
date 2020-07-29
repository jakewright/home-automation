package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jakewright/home-automation/tools/bolt/pkg/compose"
	"github.com/jakewright/home-automation/tools/bolt/pkg/service"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

var (
	stopCmd = &cobra.Command{
		Use:   "stop [foo] [bar]...",
		Short: "stop a service",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := compose.New()
			if err != nil {
				output.Fatal("Failed to init compose: %v", err)
			}

			all, err := cmd.Flags().GetBool("all")
			if err != nil {
				output.Fatal("Failed to parse all flag: %v", err)
			}

			if all {
				if err := c.StopAll(); err != nil {
					output.Fatal("Failed to stop services: %v", err)
				}

				return
			}

			if len(args) == 0 {
				return
			}

			services := service.Expand(args)

			if err := c.Stop(services); err != nil {
				output.Fatal("Failed to stop services: %v", err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().Bool("all", false, "stop all services")
}
