package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/vogtp/go-system-check/cmd/root"
)

var (
	kb float64 = 1024
	mb         = kb * kb
	gb         = mb * kb
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	root.Command(ctx)
}
func main2() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	v, err := mem.VirtualMemory()
	if err != nil {
		slog.Warn("Cannot stat virtal memory", "err", err)
	}

	// almost every return value is a struct
	fmt.Printf("VMem Total: %v, Free:%v, UsedPercent:%.2f%%\n", v.Total, v.Free, v.UsedPercent)

	parts, err := disk.PartitionsWithContext(ctx, true)
	if err != nil {
		slog.Warn("Cannot stat parts", "err", err)
	}
	for _, p := range parts {
		du, err := disk.UsageWithContext(ctx, p.Mountpoint)
		if err != nil {
			slog.Warn("Cannot get partition usage", "mountpoint", p.Mountpoint)
		}
		fmt.Printf("%s Used: %.2f%%  %.2f/%.2f %.2f\n", p.Mountpoint, du.UsedPercent, float64(du.Used)/gb, float64(du.Total)/gb, float64(du.Free)/gb)
	}

	/*
		for range 5 {
			cpuLoad, err := load.AvgWithContext(ctx)
			if err != nil {
				slog.Warn("Cannot stat cpu load", "err", err)
				break
			}
			if cpuLoad.Load1 == 0 {
				if ctx.Err() != nil {
					break
				}
				slog.Info("Load is 0", "cpuLoad", cpuLoad, "err", err)
				time.Sleep(6 * time.Second)
				continue
			}
			fmt.Printf("%v\n", cpuLoad)
		}
	*/
}
