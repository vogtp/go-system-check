package net

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/check"
	"github.com/vogtp/go-icinga/pkg/icinga"
	"github.com/vogtp/go-system-check/pkg/net"
)

const (
	ignoredProtos = "netstat.ignore-proto"
	detailsCnt    = "netstat.conn_cnt"
)

// Command adds all memory commands
func Command() *cobra.Command {
	netCmd.AddCommand(netStatCmd)
	flags := netCmd.PersistentFlags()
	flags.StringSlice(ignoredProtos, []string{"unix"}, "Protocolls to be ignored")
	flags.Int(detailsCnt, 1, "Minimum number of connections to a remote system, to be listed in the details")
	flags.VisitAll(func(f *pflag.Flag) {
		if err := viper.BindPFlag(f.Name, f); err != nil {
			panic(err)
		}
	})
	return netCmd
}

var netCmd = &cobra.Command{
	Use:   "net",
	Short: "Monitor network parameters",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		return cmd.Help()
	},
}

var netStatCmd = &cobra.Command{
	Use:     "stat",
	Short:   "Monitor network connections",
	Aliases: []string{"conn"},
	Long:    ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		check.SetWarningThresholdDefault("10000")
		check.SetCriticalThresholdDefault("20000")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		result := check.NewResult(cmd.CommandPath())
		defer result.PrintExit()
		ns, err := net.Stat(viper.GetStringSlice(ignoredProtos))
		if err != nil {
			result.SetCode(icinga.UNKNOWN)
			slog.Warn("Cannot run netstat", "err", err)
			return err
		}
		total := 0
		for k, v := range ns.Summary {
			result.SetCounter(k, v)
			total += v
		}
		result.SetCounter("total", total)
		minCnt := viper.GetInt(detailsCnt)
		for s, det := range ns.Deatils {
			statusHeader := false
			for r, cnt := range det {
				if cnt <= minCnt {
					continue
				}
				if !statusHeader {
					result.SetStatus(s, "count")
					statusHeader = true
				}
				result.SetStatus(r, cnt)
			}
		}
		result.SetHeader("%s %v", "total connections:", total)
		return nil
	},
}
