package cmd

import "github.com/spf13/cobra"

var (
	schemaCmd = &cobra.Command{}
)

func init() {
	rootCmd.AddCommand(schemaCmd)
}
