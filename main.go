package main

import (
	"fmt"

	goicinga "github.com/vogtp/go-icinga"
	"github.com/vogtp/go-icinga/pkg/check"
	"github.com/vogtp/go-icinga/pkg/icinga"
	"github.com/vogtp/go-system-check/cmd/cpu"
	"github.com/vogtp/go-system-check/cmd/disk"
	"github.com/vogtp/go-system-check/cmd/hashcmd"
	"github.com/vogtp/go-system-check/cmd/memory"
	"github.com/vogtp/go-system-check/cmd/systemdcmd"
)

func init() {
	goicinga.VersionMajor = 0
	goicinga.VersionMinor = 4
	goicinga.VersionPatch = 0
}

func main() {
	rootCtl := &check.Command{
		Use:             "syscheck",
		Short:           "Selfcontained icinga system checks",
		NamePrefix:      "293",
		DescriptionURL:  "https://github.com/vogtp/go-system-check",
		Criticality:     icinga.Criticality7x24,
		DefaultRemoteOn: true,
	}
	rootCtl.AddCommand(cpu.Command())
	// rootCtl.AddCommand(testcmd.Command())
	rootCtl.AddCommand(hashcmd.Command())
	rootCtl.AddCommand(memory.Command())
	rootCtl.AddCommand(disk.Command())
	rootCtl.AddCommand(systemdcmd.Command())

	if err := rootCtl.Execute(); err != nil {
		fmt.Println(err)
	}
}
