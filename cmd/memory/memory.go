package memory

import (
	"fmt"
	"log/slog"

	"github.com/shirou/gopsutil/v4/mem"
	"github.com/spf13/cobra"
	"github.com/vogtp/go-icinga/pkg/check"
	"github.com/vogtp/go-icinga/pkg/icinga"
	"github.com/vogtp/go-icinga/pkg/unit"
)

// Command adds all memory commands
func Command() *cobra.Command {
	return memoryCmd
}

const (
	usedPercent = "used_percent"
)

func memoryFormater() check.CheckResultOption {
	return check.CounterFormater(func(name string, value check.Value) string {
		f, ok := value.Value.(float64)
		if !ok {
			return unit.FormatGB(value.Value)
		}
		return fmt.Sprintf("%.3f%%", f)
	},
	)
}

var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Show memory",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		check.SetWarningThresholdDefault("90%")
		check.SetCriticalThresholdDefault("98%")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		result := check.NewResult(cmd.CommandPath(), memoryFormater())

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
		result.SetHeader("Used %.0f%%", v.UsedPercent)
		return nil
	},
}
