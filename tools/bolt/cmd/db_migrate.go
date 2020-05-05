package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jakewright/home-automation/tools/bolt/pkg/config"
	"github.com/jakewright/home-automation/tools/bolt/pkg/database"
	"github.com/jakewright/home-automation/tools/bolt/pkg/service"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

var (
	dbSchemaCmd = &cobra.Command{
		Use:   "schema [service.foo] [service.bar]",
		Short: "apply a service's schema",
		Long:  "Apply one or more services' schemas. If run without arguments, all service schemas are applied.",
		Run: func(cmd *cobra.Command, args []string) {
			// If no arguments supplied, assume all services.
			if len(args) == 0 {
				var err error
				args, err = service.ListAll()
				if err != nil {
					output.Fatal("Failed to list all services: %v", err)
				}
			}

			db := config.Get().Database

			for _, serviceName := range args {
				schema, err := database.GetDefaultSchema(serviceName)
				if err != nil {
					output.Fatal("Failed to get schema for %s: %v", serviceName, err)
				}

				// Silently skip services that don't have a schema
				if schema == "" {
					continue
				}

				if err := database.ApplySchema(&db, serviceName, schema); err != nil {
					output.Fatal("Failed to apply schema: %v", err)
				}
			}
		},
	}
)

func init() {
	dbCmd.AddCommand(dbSchemaCmd)
}
