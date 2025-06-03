package memory

import (
	"fmt"

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

		v, err := mem.VirtualMemoryWithContext(ctx)
		if err != nil {
			return err
		}
		result.Counter["total"] = v.Total
		result.Counter["used"] = v.Used
		result.Counter["free"] = v.Free
		result.Counter[usedPercent] = v.UsedPercent

		result.PrintExit()
		return nil
	},
}
