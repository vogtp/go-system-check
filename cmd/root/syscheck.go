package root

import (
	"context"
	"fmt"
	"log/slog"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/director"
	cpucmd "github.com/vogtp/go-system-check/cmd/cpu"
	"github.com/vogtp/go-system-check/cmd/disk"
	"github.com/vogtp/go-system-check/cmd/hashcmd"
	"github.com/vogtp/go-system-check/cmd/memory"
	"github.com/vogtp/go-system-check/cmd/systemdcmd"
	"github.com/vogtp/go-system-check/pkg/ssh"
)

// Command adds the root command
func Command(ctx context.Context) {
	rootCtl.AddCommand(cpucmd.Command())
	// rootCtl.AddCommand(testcmd.Command())
	rootCtl.AddCommand(hashcmd.Command())
	rootCtl.AddCommand(memory.Command())
	rootCtl.AddCommand(disk.Command())
	rootCtl.AddCommand(systemdcmd.Command())

	flags := rootCtl.PersistentFlags()
	ssh.Flags(flags)
	director.Flags(flags)
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
			//	fmt.Printf("ssh key: %s\n", viper.GetString("remote.sshkey"))
			if err := ssh.RemoteCheck(cmd, args); err != nil {
				u, err2 := user.Current()
				slog.Warn("Remote check error", "username", u.Name, "home", u.HomeDir, "errr", err2)
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)
