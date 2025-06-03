package hashcmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vogtp/go-system-check/pkg/hash"
)

// Command adds all hash commands
func Command() *cobra.Command {
	hashCmd.AddCommand(hashCheckCmd)
	return hashCmd
}

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Show file hash",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := hash.Calc()
		fmt.Printf("%s\n", h)
		return err
	},
}
var hashCheckCmd = &cobra.Command{
	Use:   "check <hash>",
	Short: "check file hash",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Help()
		}
	
		if err:= hash.Check(args[0]); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		return nil
	},
}
