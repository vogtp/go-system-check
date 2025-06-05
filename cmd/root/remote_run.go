package root

import (
	"fmt"
	"html"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/checks"
	"github.com/vogtp/go-system-check/pkg/ssh"
)

var ignoredFlags = []string{"help", remoteHost, remoteUser}

func remoteRun(cmd *cobra.Command, args []string) error {
	checks.InitLog()
	rh := viper.GetString(remoteHost)
	if len(rh) > 0 && rh != remoteHostDefault {
		cmds := strings.Split(cmd.CommandPath(), " ")
		cmds = append(cmds, args...)
		cmd.Flags().Visit(func(f *pflag.Flag) {
			if slices.Contains(ignoredFlags, f.Name) {
				return
			}
			val := f.Value.String()
			if strings.HasSuffix(f.Value.Type(), "Slice") {
				val = strings.ReplaceAll(val, "[", "")
				val = strings.ReplaceAll(val, "]", "")
				val = strings.ReplaceAll(val, ", ", ",")
			}
			cmds = append(cmds, fmt.Sprintf("--%s", f.Name), val)
		})
		slog.Debug("Handle remote run", remoteHost, rh, remoteUser, viper.GetString(remoteUser), "commands", cmd.CommandPath(), "args", cmd.Args, "command", cmds)

		out, err := ssh.RunOrCopy(cmd.Context(), viper.GetString(remoteUser), viper.GetString(remoteHost), cmds)
		if err != nil {
			return err
		}
		if checks.LogBuffer.Len() > 0 {
			out = strings.ReplaceAll(out, "|", fmt.Sprintf("\nLocal Log:\n%s|", html.EscapeString(checks.LogBuffer.String())))
		}
		if len(out) < 1 {
			fmt.Println(checks.LogBuffer.String())
			os.Exit(1)
		}
		fmt.Print(out)
		os.Exit(0)
	}
	return nil
}
