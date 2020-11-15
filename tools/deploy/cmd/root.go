package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/jakewright/home-automation/tools/deploy/pkg/config"
	"github.com/jakewright/home-automation/tools/deploy/pkg/deployer"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
	"github.com/jakewright/home-automation/tools/libraries/cache"
)

var (
	configPath = "./private/deploy/config.yml"

	rootCmd = &cobra.Command{
		Use:       "deploy",
		Short:     "A deployment tool for home automation",
		ValidArgs: []string{"service"},
		Args:      cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				output.Fatal("Failed to get verbose flag: %v", err)
			}
			output.Verbose = verbose

			if err := config.Init(configPath); err != nil {
				output.Fatal("Failed to load config: %v", err)
			}

			service := config.Get().Services[args[0]]
			if service == nil {
				output.Fatal("Unknown service %q", args[0])
				return
			}

			if len(service.Targets()) == 0 {
				output.Fatal("Service has no targets")
			}

			target := service.Targets()[0]

			if len(service.Targets()) > 1 {
				prompt := promptui.Select{
					Label: "Select target",
					Items: service.TargetNames(),
				}

				if i, _, err := prompt.Run(); err != nil {
					output.Fatal("Prompt failed: %v", err)
				} else {
					target = service.Targets()[i]
				}
			}

			deployer, err := deployer.Choose(service, target)
			if err != nil {
				output.Fatal("Failed to choose deployment method: %v", err)
			}

			getRevision, err := cmd.Flags().GetBool("revision")
			if err != nil {
				output.Fatal("Failed to get revision flag: %v", err)
			}

			if getRevision {
				revision, err := deployer.Revision()
				if err != nil {
					output.Fatal("Failed to get revision: %v", err)
				}

				if revision == "" {
					output.Info("\nNot yet deployed\n")
					return
				}

				output.Info("\nCurrently deployed revision: %s\n", revision)
				return
			}

			revision := "master"
			if len(args) > 1 {
				revision = args[1]
			}

			if err := deployer.Deploy(revision); err != nil {
				output.Fatal("Failed to deploy: %v", err)
			}
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
	if err := cache.Init("deploy"); err != nil {
		output.Fatal("Failed to initialise cache: %v", err)
	}

	rootCmd.PersistentFlags().Bool("revision", false, "Retrieve the currently deployed version of the service")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Display verbose output")
}
