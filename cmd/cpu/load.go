package cpu

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/check"
	"github.com/vogtp/go-icinga/pkg/icinga"
)

const (
	proocListCnt = "process.list.count"
)

var cpuLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Show cpu load",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		check.SetCriticalThresholdDefault("total:99%")
		check.SetWarningThresholdDefault("total:90%")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		result := check.NewResult(cmd.CommandPath(), check.PercentCounterFormater())

		defer result.PrintExit()
		cpuPercent, err := cpu.PercentWithContext(ctx, 200*time.Millisecond, true)
		if err != nil {
			slog.Warn("Cannot stat cpu percent", "err", err)
			result.SetCode(icinga.UNKNOWN)
			return err
		}
		var t float64
		for i, f := range cpuPercent {
			result.SetCounter(fmt.Sprintf("cpu%v", i), f)
			// fmt.Printf("cpu%v %.3f%%\n", i, f)
			t += f
		}
		total := t / float64(len(cpuPercent))
		result.SetHeader("Total load %v", total)
		result.SetCounter("total", total)
		// fmt.Printf("total %.3f%%\n", t/float64(len(cpuPercent)))

		if err := listTopProcesses(ctx, result); err != nil {
			return err
		}

		return nil
	},
}

func listTopProcesses(ctx context.Context, result *check.Result) error {
	procs, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return fmt.Errorf("cannot list processes: %w", err)
	}
	slices.SortFunc(procs, func(a, b *process.Process) int {
		aP, err := a.CPUPercent()
		if err != nil {
			slog.Warn("Cannot get process CPU%", "err", err)
		}
		bP, err := b.CPUPercent()
		if err != nil {
			slog.Warn("Cannot get process CPU%", "err", err)
		}
		return cmp.Compare(bP, aP)
	})
	for i, p := range procs {
		if i >= viper.GetInt(proocListCnt) {
			break
		}
		n, err := p.Name()
		if err != nil {
			slog.Warn("Cannot get process name", "err", err)
		}
		cpuPer, err := p.CPUPercent()
		if err != nil {
			slog.Warn("Cannot get process CPU%", "err", err)
		}
		result.SetStatus(n, fmt.Sprintf("%.1f%%", cpuPer))
	}
	return nil
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
			cpuPercent, _ := cpu.PercentWithContext(ctx, 200*time.Millisecond, true)
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
