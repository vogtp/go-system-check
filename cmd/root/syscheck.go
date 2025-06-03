package root

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/director"
	cpucmd "github.com/vogtp/go-system-check/cmd/cpu"
	"github.com/vogtp/go-system-check/cmd/hashcmd"
	"github.com/vogtp/go-system-check/cmd/memory"
)

const (
	remoteHost        = "remote.host"
	remoteUser        = "remote.user"
	remoteHostDefault = "$host.name$"
)

// Command adds the root command
func Command(ctx context.Context) {
	rootCtl.AddCommand(cpucmd.Command())
	// rootCtl.AddCommand(testcmd.Command())
	rootCtl.AddCommand(hashcmd.Command())
	rootCtl.AddCommand(memory.Command())

	flags := rootCtl.PersistentFlags()
	flags.String(remoteHost, remoteHostDefault, "Remote host to run the command on")
	flags.String(remoteUser, "root", "Remote user name")
	director.GenerateDirectorConfigPFlag(flags)
	flags.VisitAll(func(f *pflag.Flag) {
		if err := viper.BindPFlag(f.Name, f); err != nil {
			panic(err)
		}
	})

	if err := rootCtl.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
	}
}

var (
	rootCtl = &cobra.Command{
		Use:   "syscheck",
		Short: "Selfcontained icinga system checks",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := generateDirectorConfig(cmd, args); err != nil {
				return err
			}
			if err := remoteRun(cmd, args); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)
