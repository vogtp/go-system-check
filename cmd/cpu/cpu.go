package cpu

import (
	"fmt"
	"log/slog"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Command adds all cpu commands
func Command() *cobra.Command {
	cpuCmd.AddCommand(cpuListCmd)
	cpuCmd.AddCommand(cpuLoadCmd)
	flags := cpuCmd.PersistentFlags()
	flags.Int(proocListCnt, 5, "Number of top processes to be listed")
	flags.VisitAll(func(f *pflag.Flag) {
		if err := viper.BindPFlag(f.Name, f); err != nil {
			panic(err)
		}
	})
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
