package cmd

import "github.com/spf13/cobra"

var (
	dbCmd = &cobra.Command{
		Use:   "db [command]",
		Short: "perform database operations",
	}
)

func init() {
	rootCmd.AddCommand(dbCmd)
}
