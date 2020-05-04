package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jakewright/home-automation/tools/bolt/pkg/config"
	"github.com/jakewright/home-automation/tools/bolt/pkg/service"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
	"github.com/jakewright/home-automation/tools/toolutils"
)

var (
	rootCmd = &cobra.Command{
		Use:   "run service.foo",
		Short: "A tool to run home automation services locally",
		Long:  "Long version",
		Args:  cobra.MinimumNArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Global setup goes here
			if err := toolutils.Init("run"); err != nil {
				output.Fatal("Failed to initialise toolutils: %v", err)
			}

			if err := config.Init(); err != nil {
				output.Fatal("Failed to initialise config: %v", err)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			build, err := cmd.Flags().GetBool("build")
			if err != nil {
				output.Fatal("Failed to parse build flag: %v", err)
			}

			if build {
				if err := service.Build(args); err != nil {
					output.Fatal("Failed to build: %v", err)
				}
			}

			if err := service.Run(args); err != nil {
				output.Fatal("Failed to run: %v", err)
			}
		},
	}
)

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("build", "b", false, "rebuild the service before running")
}
