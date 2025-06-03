package root

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vogtp/go-system-check/pkg/ssh"
)

func remoteRun(cmd *cobra.Command, args []string) error {
	rh := viper.GetString(remoteHost)
	if len(rh) > 0 && rh != remoteHostDefault {
		cmds := strings.Split(cmd.CommandPath(), " ")
		slog.Info("Handle remote run", remoteHost, rh, remoteUser, viper.GetString(remoteUser), "commands", cmd.CommandPath(), "args", os.Args)
		cmds = append(cmds, args...)
		err := ssh.RunOrCopy(cmd.Context(), viper.GetString(remoteUser), viper.GetString(remoteHost), cmds)
		if err != nil {
			return err
		}
		os.Exit(0)
	}
	return nil
}
