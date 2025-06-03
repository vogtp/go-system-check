package systemdcmd

import (
	"fmt"

	"github.com/shirou/gopsutil/mem"
	"github.com/spf13/cobra"
	"github.com/vogtp/go-icinga/pkg/checks"
	"github.com/vogtp/go-icinga/pkg/icinga"
)

// Command adds all memory commands
func Command() *cobra.Command {
	systemdCmd.AddCommand(systemdUnitCmd)
	return systemdCmd
}

var systemdCmd = &cobra.Command{
	Use:   "systemd",
	Short: "Monitor systemd",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		return cmd.Help()
	},
}

var systemdUnitCmd = &cobra.Command{
	Use:   "unit",
	Short: "check a systemd unit",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		result := checks.Result{
			Name:    cmd.CommandPath(),
			Prefix:  "",
			Result:  icinga.OK,
			Stati:   make(map[string]any),
			Counter: make(map[string]any),
			CounterFormater: func(name string, value any) string {
				f, ok := value.(float64)
				if !ok {
					return fmt.Sprintf("%v", value)
				}
				return fmt.Sprintf("%.3f%%", f)
			},
		}

		

		result.PrintExit()
		return nil
	},
}
