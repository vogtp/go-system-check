package ssh

import (
	"context"
	"fmt"
	"os"

	scp "github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

func Copy(ctx context.Context, sshClient *ssh.Client, local, remote string) error {
	client, err := scp.NewClientBySSH(sshClient)
	if err != nil {
		return fmt.Errorf("error creating new SSH session from existing connection: %w", err)
	}
	defer client.Close()

	f, err := os.Open(local)
	if err != nil {
		return fmt.Errorf("error opening local file %q: %w", local, err)
	}
	defer f.Close()

	if err := client.CopyFromFile(ctx, *f, remote, "0755"); err != nil {
		return fmt.Errorf("error while scp file %q: %w", remote, err)
	}
	return nil
}
