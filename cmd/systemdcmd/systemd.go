package systemdcmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/check"
	"github.com/vogtp/go-icinga/pkg/icinga"
	"github.com/vogtp/go-system-check/pkg/systemd"
)

const (
	systemdUnits = "systemd.services"
)

// Command adds all memory commands
func Command() *cobra.Command {
	systemdCmd.AddCommand(systemdServiceCmd)
	flags := systemdCmd.PersistentFlags()
	flags.StringSlice(systemdUnits, []string{"ssh"}, "Systemd services to be checked (e.g. ssh,ufw,httpd) ")
	flags.VisitAll(func(f *pflag.Flag) {
		if err := viper.BindPFlag(f.Name, f); err != nil {
			panic(err)
		}
	})
	return systemdCmd
}

var systemdCmd = &cobra.Command{
	Use:   "systemd",
	Short: "Monitor systemd",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		return cmd.Help()
	},
}

var systemdServiceCmd = &cobra.Command{
	Use:     "service <service> ...",
	Short:   "check a systemd unit",
	Long:    ``,
	Aliases: []string{"unit"},
	RunE: func(cmd *cobra.Command, args []string) error {
		units := append(args, viper.GetStringSlice(systemdUnits)...)
		if len(units) < 1 {
			slog.Warn("No services given", "args", args, "systemdUnits", viper.GetString(systemdUnits), "units", units)
			return cmd.Help()
		}
		result := check.NewResult(cmd.CommandPath(), check.CounterFormater(activeStateFormater), check.DisplayFormater(systemdUnitTableFormater))

		defer result.PrintExit()
		var h strings.Builder
		for _, unit := range units {
			service, err := systemd.Unit(unit)
			if err != nil {
				result.SetCode(icinga.UNKNOWN)
				slog.Warn("Cannot check systemd", "unit", unit, "err", err)
				return err
			}
			result.SetCounter(unit, service)
			code := getResultCode(service)
			result.SetCode(code)
			h.WriteString(fmt.Sprintf("%s [%s] ", unit, code))
		}
		result.SetHeader("%s", h.String())
		return nil
	},
}

func getResultCode(service *systemd.Service) icinga.ResultCode {
	if service.ActiveStateInt() < 1 {
		return icinga.WARNING
	}
	if service.ActiveStateInt() < 0 && service.Preset() == "enabled" {
		return icinga.CRITICAL
	}
	return icinga.OK
}

func activeStateFormater(name string, value any) string {
	f, ok := value.(*systemd.Service)
	if !ok {
		return fmt.Sprintf("%T", value)
	}
	return fmt.Sprintf("%v", f.ActiveStateInt())
}

func systemdUnitTableFormater(counter map[string]any) string {
	rowHeader := table.Row{"Unit", "State", "Preset"}
	rows := make([]table.Row, 0, len(counter))
	for n, v := range counter {
		u, ok := v.(*systemd.Service)
		if !ok {
			slog.Warn("Not a systemd.Service", "counter", v)
			continue
		}
		rows = append(rows, table.Row{n, u.ActiveState(), u.Preset()})
	}

	//slices.SortFunc(rows, tableSort)
	tw := table.NewWriter()
	tw.AppendHeader(rowHeader)
	tw.AppendRows(rows)
	tw.SetIndexColumn(0)
	style := table.StyleLight
	style.HTML.EscapeText = true
	tw.SetStyle(style)
	// tw.SetColumnConfigs([]table.ColumnConfig{
	// 	// 	//{Name: colTitleFirstName, Align: text.AlignRight},
	// 	// 	// the 5th column does not have a title, so use the column number as the
	// 	// 	// identifier for the column
	// 	// 	{Number: 0, WidthMax: 1},
	// 	{Name: "Status", WidthMin: len("CRITICAL")},
	// })
	return tw.Render()
}

// func tableSort(a, b table.Row) int {
// 	if len(a) < 1 || len(b) < 1 {
// 		return 0
// 	}
// 	return len(a[0].(string)) - len(b[0].(string))
// }
