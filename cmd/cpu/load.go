package cpu

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/spf13/cobra"
)

var cpuLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Show cpu load",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cpuPercent, err := cpu.PercentWithContext(ctx, 200*time.Millisecond, true)
		if err != nil {
			slog.Warn("Cannot stat cpu percent", "err", err)
			return err
		}
		var t float64
		for i, c := range cpuPercent {
			fmt.Printf("cpu%v %.3f%%\n", i, c)
			t += c
		}
		fmt.Printf("total %.3f%%\n", t/float64(len(cpuPercent)))

		return nil
	},
}

var cpuLoadFollowCmd = &cobra.Command{
	Use:     "follow",
	Short:   "Show cpu load",
	Long:    ``,
	Aliases: []string{"f"},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		tick := time.NewTicker(1 * time.Second).C
		for range 500 {
			cpuPercent, err := cpu.PercentWithContext(ctx, 200*time.Millisecond, true)
			cpuPercentTot, err := cpu.PercentWithContext(ctx, 200*time.Millisecond, false)
			if err != nil {
				slog.Warn("Cannot stat cpu percent", "err", err)
				return err
			}
			var t float64
			for _, c := range cpuPercent {
				//fmt.Printf("cpu%v %.3f%%\n", i, c)
				t += c
			}
			fmt.Printf("cpu%v %.3f%% %.3f%%\n", " total", t/float64(len(cpuPercent)), cpuPercentTot[0])
			if ctx.Err() != nil {
				break
			}
			<-tick
		}
		return nil
	},
}
