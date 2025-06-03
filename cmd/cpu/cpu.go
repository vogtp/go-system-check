package cpu

import (
	"fmt"
	"log/slog"

	"github.com/shirou/gopsutil/cpu"
	"github.com/spf13/cobra"
)

// Command adds all cpu commands
func Command() *cobra.Command {
	cpuCmd.AddCommand(cpuListCmd)
	cpuCmd.AddCommand(cpuLoadCmd)
	cpuLoadCmd.AddCommand(cpuLoadFollowCmd)
	return cpuCmd
}

var cpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "Show cpu load",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var cpuListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show a list of CPUs",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cpus, err := cpu.InfoWithContext(ctx)
		if err != nil {
			slog.Warn("Cannot stat cpu info", "err", err)
			return err
		}
		for _, c := range cpus {
			fmt.Printf("cpu%v %v %v\n", c.CPU, c.Mhz, c.ModelName)
		}

		return nil
	},
}
