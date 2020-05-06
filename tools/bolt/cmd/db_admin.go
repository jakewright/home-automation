package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jakewright/home-automation/libraries/go/exe"
	"github.com/jakewright/home-automation/tools/bolt/pkg/config"
	"github.com/jakewright/home-automation/tools/bolt/pkg/service"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

var (
	dbAdminCmd = &cobra.Command{
		Use:   "admin",
		Short: "open the database admin UI",
		Run: func(cmd *cobra.Command, args []string) {
			serviceName := config.Get().Database.AdminService
			if err := service.Run([]string{serviceName}); err != nil {
				output.Fatal("Failed to start admin service %q: %v", serviceName, err)
			}

			ports, err := service.Ports(serviceName)
			if err != nil {
				output.Fatal("Failed to get admin service ports: %v", err)
			}

			if len(ports) != 1 {
				output.Fatal("Unexpected number of ports %d: %+v", len(ports), ports)
			}

			if err := exe.Command("open", "http://localhost:"+ports[0]).Run().Err; err != nil {
				output.Fatal("Failed to open URL: %v", err)
			}
		},
	}
)

func init() {
	dbCmd.AddCommand(dbAdminCmd)
}
