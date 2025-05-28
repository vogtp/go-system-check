package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	v, err := mem.VirtualMemory()
	if err != nil {
		slog.Warn("Cannot stat virtal memory", "err", err)
	}

	// almost every return value is a struct
	fmt.Printf("VMem Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	parts, err := disk.PartitionsWithContext(ctx, true)
	if err != nil {
		slog.Warn("Cannot stat parts", "err", err)
	}
	for _, p := range parts {
		du, err := disk.UsageWithContext(ctx, p.Mountpoint)
		if err != nil {
			slog.Warn("Cannot get partition usage", "mountpoint", p.Mountpoint)
		}
		fmt.Printf("%v Used: %v%%  %v/%v\n", p.Mountpoint, du.UsedPercent, du.Used, du.Total)
	}

	cpus, err := cpu.InfoWithContext(ctx)
	if err != nil {
		slog.Warn("Cannot stat cpu info", "err", err)
	}

	for _, c := range cpus {
		fmt.Printf("cpu%v %v %v\n", c.CPU, c.Mhz, c.ModelName)
	}
	for {
		cpuLoad, err := load.AvgWithContext(ctx)
		if err != nil {
			slog.Warn("Cannot stat cpu load", "err", err)
			break
		}
		if cpuLoad.Load1 == 0 {
			if ctx.Err() != nil {
				break
			}
			slog.Info("Load is 0", "cpuLoad", cpuLoad)
			time.Sleep(6 * time.Second)
			continue
		}
		fmt.Printf("%v\n", cpuLoad)
	}
}
