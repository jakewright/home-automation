package cmd

import "github.com/spf13/cobra"

var (
	firehoseCmd = &cobra.Command{
		Use:   "firehose [command]",
		Short: "interact with the Firehose",
	}
)

func init() {
	rootCmd.AddCommand(firehoseCmd)
}
