package memory

import (
	"log/slog"

	"github.com/shirou/gopsutil/mem"
	"github.com/spf13/cobra"
	"github.com/vogtp/go-icinga/pkg/checks"
	"github.com/vogtp/go-icinga/pkg/icinga"
)

// Command adds all memory commands
func Command() *cobra.Command {
	return memoryCmd
}

const (
	usedPercent = "used_percent"
)

var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Show memory",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		result := checks.NewCheckResult(cmd.CommandPath(), checks.PercentCounterFormater())

		defer result.PrintExit()

		v, err := mem.VirtualMemoryWithContext(ctx)
		if err != nil {
			result.SetCode(icinga.WARNING)
			slog.Warn("Cannot get memory", "err", err)
			return err
		}
		result.SetCounter("total", v.Total)
		result.SetCounter("used", v.Used)
		result.SetCounter("free", v.Free)
		result.SetCounter(usedPercent, v.UsedPercent)
		if v.UsedPercent > 90 {
			result.SetCode(icinga.WARNING)
		}
		if v.UsedPercent > 98 {
			result.SetCode(icinga.CRITICAL)
		}
		return nil
	},
}
