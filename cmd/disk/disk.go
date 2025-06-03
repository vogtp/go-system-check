package disk

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shirou/gopsutil/disk"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/checks"
	"github.com/vogtp/go-icinga/pkg/icinga"
)

// Command adds all memory commands
func Command() *cobra.Command {
	flags := diskCmd.PersistentFlags()
	flags.StringSlice(excludeParts, []string{"/run", "/snap", "/sys", "/dev", "/proc"}, "Partions to be excluded")
	flags.VisitAll(func(f *pflag.Flag) {
		if err := viper.BindPFlag(f.Name, f); err != nil {
			panic(err)
		}
	})
	return diskCmd
}

const (
	excludeParts = "exclude"
	usedPercent  = "used_percent"
)

var (
	kb float64 = 1024
	mb         = kb * kb
	gb         = mb * kb
)

var diskCmd = &cobra.Command{
	Use:   "disk",
	Short: "Show disk usage",
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
			DisplayCounterFormater: diskTableFormater,
		}

		parts, err := disk.PartitionsWithContext(ctx, true)
		if err != nil {
			return err
		}
		for _, p := range parts {
			if exclude(p.Mountpoint, viper.GetStringSlice(excludeParts)...) {
				continue
			}
			du, err := disk.UsageWithContext(ctx, p.Mountpoint)
			if err != nil {
				slog.Warn("Cannot get partition usage", "mountpoint", p.Mountpoint)
				continue
			}
			result.Counter[p.Mountpoint+"-total"] = du.Total
			result.Counter[p.Mountpoint+"-percent"] = du.UsedPercent
			result.Counter[p.Mountpoint+"-usage"] = du.Used
			result.Counter[p.Mountpoint+"-free"] = du.Free
		}

		result.PrintExit()
		return nil
	},
}

func exclude(path string, excl ...string) bool {
	for _, e := range excl {
		if strings.HasPrefix(path, e) {
			return true
		}
	}
	return false
}

func diskTableFormater(counter map[string]any) string {
	rowHeader := table.Row{"Partiton", "Percent", "Used", "Free", "Total"}
	disks := make(map[string]table.Row)
	for n, v := range counter {
		split := strings.Split(n, "-")
		diskName := split[0]
		d, ok := disks[diskName]
		if !ok {
			d = make([]any, 5)
			d[0] = diskName
		}
		switch split[1] {
		case "percent":
			d[1] = fmt.Sprintf("%.1f%%", v)
		case "usage":
			d[2] = formatGB(v)
		case "free":
			d[3] = formatGB(v)
		case "total":
			d[4] = formatGB(v)
		}
		disks[diskName] = d
	}
	diskRows := make([]table.Row, 0, len(disks))
	for _, d := range disks {
		diskRows = append(diskRows, d)
	}
	slices.SortFunc(diskRows, tableSort)
	tw := table.NewWriter()
	tw.AppendHeader(rowHeader)
	tw.AppendRows(diskRows)
	tw.SetIndexColumn(0)
	style := table.StyleLight
	style.HTML.EscapeText = true
	tw.SetStyle(style)
	return tw.Render()
}
func tableSort(a, b table.Row) int {
	if len(a) < 1 || len(b) < 1 {
		return 0
	}
	return len(a[0].(string)) - len(b[0].(string))
}

func formatGB(d any) string {
	i, ok := d.(uint64)
	if !ok {
		return fmt.Sprintf("%v", d)
	}
	f := float64(i)
	if f > gb {
		return fmt.Sprintf("%.0f GB", f/gb)
	}
	if f > mb {
		return fmt.Sprintf("%.0f MB", f/mb)
	}
	return fmt.Sprintf("%.0f KB", f/kb)
}
