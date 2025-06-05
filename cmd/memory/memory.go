package memory

import (
	"fmt"
	"log/slog"

	"github.com/shirou/gopsutil/mem"
	"github.com/spf13/cobra"
	"github.com/vogtp/go-icinga/pkg/checks"
	"github.com/vogtp/go-icinga/pkg/icinga"
	"github.com/vogtp/go-system-check/pkg/unit"
)

// Command adds all memory commands
func Command() *cobra.Command {
	return memoryCmd
}

const (
	usedPercent = "used_percent"
)

func memoryFormater() checks.CheckResultOption {
	return checks.CounterFormater(func(name string, value any) string {
		f, ok := value.(float64)
		if !ok {
			return unit.FormatGB(value)
		}
		return fmt.Sprintf("%.3f%%", f)
	},
	)
}

var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Show memory",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		result := checks.NewCheckResult(cmd.CommandPath(), memoryFormater())

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
		result.SetHeader("Used %.0f%%", v.UsedPercent)
		return nil
	},
}
