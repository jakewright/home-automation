package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jakewright/home-automation/tools/bolt/pkg/compose"
	"github.com/jakewright/home-automation/tools/bolt/pkg/service"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

var (
	logsCmd = &cobra.Command{
		Use:   "logs [foo] [bar]",
		Short: "show logs for a set of services (default: all services)",
		Run: func(cmd *cobra.Command, args []string) {
			services := service.Expand(args)

			c, err := compose.New()
			if err != nil {
				output.Fatal("Failed to init compose: %v", err)
			}

			if err := c.Logs(services); err != nil {
				output.Fatal("Failed to output logs: %v", err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(logsCmd)
}
