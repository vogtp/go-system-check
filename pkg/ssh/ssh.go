package ssh

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/vogtp/go-system-check/pkg/hash"
	"golang.org/x/crypto/ssh"
)

func RunOrCopy(ctx context.Context, user, host string, cmd []string) error {
	if len(cmd) < 1 {
		return fmt.Errorf("no command given: %v", cmd)
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(ssh_key, ssh_key_pass)
	if err != nil {
		return fmt.Errorf("unable to parse private key: %w", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), config)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()

	h, err := hash.Calc()
	if err != nil {
		return fmt.Errorf("cannot calculate my hash: %w", err)
	}
	remote := cmd[0]
	if err := exec(client, fmt.Sprintf("./%s hash check %s", remote, h)); err != nil {
		local := os.Args[0]
		fmt.Printf("remote version is outdated: copy %q to remote file: %q\n", local, remote)
		if err := Copy(ctx, client, local, remote); err != nil {
			return err
		}
	}
	cmdLine := fmt.Sprintf("./%s", strings.Join(cmd, " "))
	slog.Info("Executing remote command", "cmd", cmdLine)
	if err := exec(client, cmdLine); err != nil {
		return fmt.Errorf("%q returned: %w", cmdLine, err)
	}
	return nil
}

func exec(client *ssh.Client, cmd string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("failed to run command %q: %w", cmd, err)
	}
	fmt.Println(b.String())
	//	session.
	return nil
}
