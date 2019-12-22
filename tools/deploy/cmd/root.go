package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/jakewright/home-automation/tools/deploy/config"
	"github.com/spf13/cobra"
)

var (
	revision bool

	rootCmd = &cobra.Command{
		Use:       "deploy",
		Short:     "A deployment tool for home automation",
		ValidArgs: []string{"service"},
		Args:      cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if revision {
				fmt.Printf("Revision flag set\n")
			}

			service := config.FindService(args[0])
			if service == nil {
				log.Fatalf("Unknown service '%s'", args[0])
			}

			switch service.System {
			case config.SysSystemd:
				//
			}

			fmt.Printf("Building %s...\n", service)
		},
	}
)

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&revision, "revision", false, "Retrieve the currently deployed version of the service")
}
