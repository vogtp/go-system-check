package ssh

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var ignoredFlags = []string{"help", remoteHost, remoteUser, remoteUserKey, remoteUserKeyPass}

const (
	remoteHost        = "remote.host"
	remoteUser        = "remote.user"
	remoteUserKey     = "remote.sshkey"
	remoteUserKeyPass = "remote.sshkeypass"
	remoteHostDefault = "$host.name$"
)

func Flags(flags *pflag.FlagSet) {
	flags.String(remoteHost, remoteHostDefault, "Remote host to run the command on")
	flags.String(remoteUser, "root", "Remote user name")
	flags.String(remoteUserKey, "/var/lib/nagios/.ssh/icinga_ssh", "ssh private key file location")
	flags.String(remoteUserKeyPass, "", "ssh private key password")
}

func IsRemoteRun() bool {
	rh := viper.GetString(remoteHost)
	return len(rh) > 0 && rh != remoteHostDefault
}
