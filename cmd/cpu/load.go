package cpu

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/spf13/cobra"
	"github.com/vogtp/go-icinga/pkg/checks"
	"github.com/vogtp/go-icinga/pkg/icinga"
)

var cpuLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Show cpu load",
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
		cpuPercent, err := cpu.PercentWithContext(ctx, 200*time.Millisecond, true)
		if err != nil {
			slog.Warn("Cannot stat cpu percent", "err", err)
			return err
		}
		var t float64
		for i, f := range cpuPercent {
			result.Counter[fmt.Sprintf("cpu%v", i)] = f
			// fmt.Printf("cpu%v %.3f%%\n", i, f)
			t += f
		}
		result.Total = t / float64(len(cpuPercent))
		result.Counter["total"] = result.Total
		// fmt.Printf("total %.3f%%\n", t/float64(len(cpuPercent)))
		result.PrintExit()
		return nil
	},
}

var cpuLoadFollowCmd = &cobra.Command{
	Use:     "follow",
	Short:   "Show cpu load follow",
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
