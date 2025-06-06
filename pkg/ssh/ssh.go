package ssh

//go:generate go run gen.go arg1 arg2

import (
	"bytes"
	"context"
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
	"github.com/vogtp/go-system-check/pkg/hash"
	"golang.org/x/crypto/ssh"
)

func RemoteCheck(cmd *cobra.Command, args []string) error {
	checks.InitLog()
	if IsRemoteRun() {
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

		out, err := runOrCopy(cmd.Context(), cmds)
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

func runOrCopy(ctx context.Context, cmd []string) (string, error) {
	if len(cmd) < 1 {
		return "", fmt.Errorf("no command given: %v", cmd)
	}
	user := viper.GetString(remoteUser)
	host := viper.GetString(remoteHost)

	sshAuth, err := getSshAuth()
	if err != nil {
		return "", fmt.Errorf("no ssh auth: %w", err)
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            sshAuth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()

	h, err := hash.Calc()
	if err != nil {
		return "", fmt.Errorf("cannot calculate my hash: %w", err)
	}
	remote := cmd[0]
	if _, err := exec(client, fmt.Sprintf("./%s hash check %s", remote, h)); err != nil {
		local := os.Args[0]
		slog.Info("remote version is outdated: copy local to remote ", "local", local, "remote", remote)
		if err := Copy(ctx, client, local, remote); err != nil {
			return "", err
		}
	}
	cmdLine := fmt.Sprintf("./%s", strings.Join(cmd, " "))
	slog.Debug("Executing remote command", "cmd", cmdLine, "host", host, "user", user)
	out, err := exec(client, cmdLine)
	if err != nil {
		return "", fmt.Errorf("%q returned: %w", cmdLine, err)
	}
	return out, nil
}

func exec(client *ssh.Client, cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()
	var stdo bytes.Buffer
	var stde bytes.Buffer
	session.Stdout = &stdo
	session.Stderr = &stde
	err = session.Run(cmd)
	if stde.Len() > 0 {
		fmt.Println(stde.String())
	}
	return stdo.String(), err
}
