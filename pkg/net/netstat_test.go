package net

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestStat(t *testing.T) {
	tests := []struct {
		name    string
		out     io.Reader
		ignored []string
		res     map[string]int
	}{
		{name: "linux", out: readFile(t, "testfiles/netstat_lnx.out"), res: map[string]int{"unix/STREAM": 1057, "tcp/ESTABLISHED": 14, "tcp6/ESTABLISHED": 10, "udp6/ESTABLISHED": 1, "unix/DGRAM": 53, "unix/SEQPACKET": 24, "tcp/CLOSE_WAIT": 1, "udp/ESTABLISHED": 2, "Active/servers)": 1}},
		{name: "linux no unix sockets", ignored: []string{"unix"}, out: readFile(t, "testfiles/netstat_lnx.out"), res: map[string]int{"tcp/ESTABLISHED": 14, "tcp/CLOSE_WAIT": 1, "tcp6/ESTABLISHED": 10, "udp/ESTABLISHED": 2, "udp6/ESTABLISHED": 1, "Active/servers)": 1}},
		{name: "windows", out: readFile(t, "testfiles/netstat_win.out"), res: map[string]int{"TCP/ESTABLISHED": 13, "TCP/TIME_WAIT": 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.name)
			np, err := parseNetStatOut(tt.out, tt.ignored)
			if err != nil {
				t.Errorf("parsing %s got error: %v", tt.name, err)
			}
			for k, v := range tt.res {
				if np.Summary[k] != v {
					t.Errorf("%s: %s expected %v got %v", tt.name, k, v, np.Summary[k])
				}
			}
		})
	}
}

func readFile(t *testing.T, name string) io.Reader {
	f, err := os.Open(name)
	if err != nil {
		t.Fatalf("Cannot read file: %v", err)
	}
	return f
}
