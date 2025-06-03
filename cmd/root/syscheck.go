package root

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	cpucmd "github.com/vogtp/go-system-check/cmd/cpu"
	"github.com/vogtp/go-system-check/cmd/hashcmd"
	testcmd "github.com/vogtp/go-system-check/cmd/test"
)

// Command adds the root command
func Command(ctx context.Context) {
	rootCtl.AddCommand(cpucmd.Command())
	rootCtl.AddCommand(testcmd.Command())
	rootCtl.AddCommand(hashcmd.Command())
	if err := rootCtl.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
	}
}

var (
	rootCtl = &cobra.Command{
		Use:   "syscheck",
		Short: "Selfcontained icinga system checks",

		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)
