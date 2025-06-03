package root

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vogtp/go-icinga/pkg/director"
	"github.com/vogtp/go-icinga/pkg/icinga"
)

func generateDirectorConfig(cmd *cobra.Command, args []string) error {
	if director.ShouldGenerate() {
		d := director.Generator{
			NamePrefix:     "293",
			Description:    "syscheck: self coping icinga remote checks",
			DescriptionURL: "https://github.com/vogtp/go-system-check",
			CobraCmd:       cmd,
			Output:         os.Stdout,
			Criticality:    icinga.Criticality7x24,
		}
		if err := d.Generate(); err != nil {
			return err
		}
		os.Exit(0)
	}
	return nil
}
