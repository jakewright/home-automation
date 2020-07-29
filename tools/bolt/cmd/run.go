package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jakewright/home-automation/tools/bolt/pkg/compose"
	"github.com/jakewright/home-automation/tools/bolt/pkg/service"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

var (
	runCmd = &cobra.Command{
		Use:   "run [foo] [bar]...",
		Short: "run a service",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			build, err := cmd.Flags().GetBool("build")
			if err != nil {
				output.Fatal("Failed to parse build flag: %v", err)
			}

			services := service.Expand(args)

			c, err := compose.New()
			if err != nil {
				output.Fatal("Failed to init compose: %v", err)
			}

			if build {
				if err := c.Build(services); err != nil {
					output.Fatal("Failed to build: %v", err)
				}
			}

			if err := service.Run(c, services); err != nil {
				output.Fatal("Failed to run: %v", err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("build", "b", false, "rebuild the service before running")
}
