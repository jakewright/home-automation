package cmd

import (
	"github.com/spf13/cobra"
)

var (
	schemaMigrateCmd = &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			//databaseServiceName := config.DatabaseService()

		},
	}
)

func init() {
	schemaCmd.AddCommand(schemaMigrateCmd)
}
