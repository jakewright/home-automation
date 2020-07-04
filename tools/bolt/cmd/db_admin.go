package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jakewright/home-automation/libraries/go/exe"
	"github.com/jakewright/home-automation/tools/bolt/pkg/compose"
	"github.com/jakewright/home-automation/tools/bolt/pkg/config"
	"github.com/jakewright/home-automation/tools/bolt/pkg/service"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

var (
	dbAdminCmd = &cobra.Command{
		Use:   "admin",
		Short: "open the database admin UI",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := compose.New()
			if err != nil {
				output.Fatal("Failed to init compose: %v", err)
			}

			serviceName := config.Get().Database.AdminService
			servicePath := config.Get().Database.AdminServicePath

			if running, err := c.IsRunning(serviceName); err != nil {
				output.Fatal("Failed to get status of %s: %v", serviceName, err)
			} else if !running {
				if err := service.Run(c, []string{serviceName}); err != nil {
					output.Fatal("Failed to start %s: %v", serviceName, err)
				}
			}

			ports, err := c.Ports(serviceName)
			if err != nil {
				output.Fatal("Failed to get admin service ports: %v", err)
			}

			if len(ports) != 1 {
				output.Fatal("Unexpected number of ports %d: %+v", len(ports), ports)
			}

			uri := filepath.Join(fmt.Sprintf("http://localhost:%s/", ports[0]), servicePath)
			if err := exe.Command("open", uri).Run().Err; err != nil {
				output.Fatal("Failed to open URL: %v", err)
			}
		},
	}
)

func init() {
	dbCmd.AddCommand(dbAdminCmd)
}
