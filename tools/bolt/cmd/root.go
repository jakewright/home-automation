package cmd

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/jakewright/home-automation/tools/bolt/pkg/config"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
	"github.com/jakewright/home-automation/tools/libraries/cache"
)

var (
	rootCmd = &cobra.Command{
		Use:   "bolt [command]",
		Short: "A tool to run home automation services locally",
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
	if err := cache.Init("run"); err != nil {
		output.Fatal("Failed to initialise toolutils: %v", err)
	}

	if err := config.Init(); err != nil {
		output.Fatal("Failed to initialise config: %v", err)
	}

	// Append the groups to the usage info
	b := bytes.Buffer{}
	w := tabwriter.NewWriter(&b, 0, 4, 5, ' ', 0)
	for name, services := range config.Get().Groups {
		sort.Strings(services)
		if _, err := fmt.Fprintf(w, "  %s\t%s\n",
			name,
			strings.Join(services, ", "),
		); err != nil {
			panic(err)
		}
	}

	if err := w.Flush(); err != nil {
		panic(err)
	}

	groupInfo := "\nGroups:\n" + b.String()
	rootCmd.SetHelpTemplate(rootCmd.HelpTemplate() + groupInfo)
}
