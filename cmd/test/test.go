package test

import (
	"github.com/spf13/cobra"
)

// Command adds all test commands
func Command() *cobra.Command {
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
