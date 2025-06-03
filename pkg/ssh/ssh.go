package ssh

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vogtp/go-system-check/pkg/hash"
	"golang.org/x/crypto/ssh"
)

func RunOrCopy(ctx context.Context, user, host string, cmd string) error {
	// key, err := os.ReadFile("/home/vogtp/.ssh/id_ed25519")
	// if err != nil {
	// 	log.Fatalf("unable to read private key: %v", err)
	// }

	// fmt.Printf("%#v\n", key)

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKeyWithPassphrase(ssh_key, ssh_key_pass)
	if err != nil {
		return fmt.Errorf("unable to parse private key: %w", err)
	}
	//var hostKey ssh.PublicKey
	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig,
	// and provide a HostKeyCallback.
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

	local := os.Args[0]
	_, remote := filepath.Split(local)
	h, err := hash.Calc()
	if err != nil {
		return fmt.Errorf("cannot calculate my hash: %w", err)
	}
	if err := exec(client, fmt.Sprintf("./%s hash check %s", remote, h)); err != nil {
		
		fmt.Printf("remote version is outdated: copy %q to remote file: %q\n", local,  remote)
		if err := Copy(ctx, client, local, remote); err != nil {
			return err
		}
	}

	if err := exec(client, fmt.Sprintf("./%s %s", remote, cmd)); err != nil {
		return err
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
