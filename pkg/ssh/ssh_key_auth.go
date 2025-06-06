package ssh

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

func getSshAuth() ([]ssh.AuthMethod, error) {
	filename := viper.GetString(remoteUserKey)
	if len(filename) < 1 {
		return nil, fmt.Errorf("please provide a ssh key file location: %s", filename)
	}
	ssh_key, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var signer ssh.Signer
	pw := viper.GetString(remoteUserKeyPass)
	if len(pw) > 0 {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(ssh_key, []byte(pw))
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key with passwordgo : %w", err)
		}
	}
	if signer == nil {
		signer, err = ssh.ParsePrivateKey(ssh_key)
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key: %w", err)
		}
	}
	return []ssh.AuthMethod{ssh.PublicKeys(signer)}, nil
}
