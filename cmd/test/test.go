package test

import (
	"github.com/spf13/cobra"
	"github.com/vogtp/go-system-check/pkg/ssh"
)

// Command adds all test commands
func Command() *cobra.Command {
	testCmd.AddCommand(testSshCmd)
	return testCmd
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test stuff",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var testSshCmd = &cobra.Command{
	Use:  "ssh <user> <host>",
	Long: ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ssh.RunOrCopy(cmd.Context(), args[0], args[1], "cpu load")
	},
}
