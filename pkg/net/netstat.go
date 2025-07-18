package net

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"slices"
	"strings"
)

func Stat(ignoredProtos []string) (*NetStat, error) {
	netstatCmd := exec.Command("netstat")
	var stdo bytes.Buffer
	var stde bytes.Buffer
	netstatCmd.Stdout = &stdo
	netstatCmd.Stderr = &stde
	if err := netstatCmd.Run(); err != nil {
		fmt.Println(stdo.String())
		fmt.Println(stde.String())
		return nil, err
	}
	np, err := parseNetStatOut(&stdo, ignoredProtos)
	return np, err
}

func parseNetStatOut(r io.Reader, ignoredProtos []string) (*NetStat, error) {
	// ignore line: Active UNIX domain sockets (w/o servers)
	ignoredProtos = append(ignoredProtos, "Active")
	np := NetStat{ignoredProtos: ignoredProtos}
	s := bufio.NewScanner(r)
	started := false
	for s.Scan() {
		t := s.Text()
		d := strings.Fields(t)
		if len(d) < 1 {
			continue
		}
		if strings.EqualFold(protoHeader, d[0]) {
			started = true
			if err := np.header(t); err != nil {
				return nil, err
			}
			continue
		}
		if !started {
			continue
		}
		// fmt.Printf("%+v %v\n", d, len(d))
		// return nil
		if err := np.parse(d); err != nil {
			return nil, err
		}
	}
	// for k, v := range np.stats {
	// 	fmt.Printf("%q: %v, ", k, v)
	// }
	// fmt.Println()
	return &np, nil
}

type NetStat struct {
	protoIdx      int
	stateIdx      int
	remoteIdx     int
	ignoredProtos []string
	Summary       map[string]int
	Deatils       map[string]map[string]int
}

const (
	protoHeader  = "proto"
	stateHeader  = "state"
	remoteHeader = "foreign"
	established  = "ESTABLISHED"
)

func (np *NetStat) header(hdrs string) error {
	headers := strings.Fields(strings.ReplaceAll(hdrs, "Address", ""))
	np.protoIdx = -1
	np.stateIdx = -1
	for i, h := range headers {
		switch strings.ToLower(h) {
		case protoHeader:
			np.protoIdx = i
		case stateHeader:
			np.stateIdx = i
		case remoteHeader:
			np.remoteIdx = i
			// default:
			// 	fmt.Printf("Unused header: %s\n", h)
		}
		if np.protoIdx > -1 && np.stateIdx > -1 {
			headers[np.protoIdx] = fmt.Sprintf("*%s*", headers[np.protoIdx])
			headers[np.stateIdx] = fmt.Sprintf("#%s#", headers[np.stateIdx])
			//fmt.Printf("%+v %v %v %v\n", headers, np.protoIdx, np.stateIdx, len(headers))
			return nil
		}
	}
	if np.protoIdx < 0 || np.stateIdx < 0 {
		return fmt.Errorf("headers not found, proto idx: %v, state ids: %v", np.protoIdx, np.stateIdx)
	}
	return nil
}

func (np *NetStat) parse(d []string) error {
	if np.Summary == nil {
		np.Summary = make(map[string]int)
		np.Deatils = make(map[string]map[string]int)
	}
	l := len(d)
	if l < np.protoIdx || l < np.stateIdx {
		return fmt.Errorf("%#v proto %v state %v len %v", d, np.protoIdx, np.stateIdx, l)
	}
	p := d[np.protoIdx]
	s := d[np.stateIdx]
	if len(p) < 1 || len(s) < 1 {
		return fmt.Errorf("netstat proto (%s) or state (%s) not found", p, s)
	}
	if slices.Contains(np.ignoredProtos, p) {
		return nil
	}
	key := fmt.Sprintf("%s/%s", p, s)
	np.Summary[key] = np.Summary[key] + 1
	// remote details
	r := d[np.remoteIdx]
	sm := np.Deatils[s]
	if sm == nil {
		sm = make(map[string]int)
		np.Deatils[s] = sm
	}
	cnt := sm[r]
	cnt++
	np.Deatils[s][r] = cnt
	return nil
}
