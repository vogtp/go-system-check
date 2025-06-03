package systemd

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Service struct {
	name string
	data map[string]string
}

func Unit(unit string) (s *Service, err error) {
	s = &Service{
		name: unit,
		data: make(map[string]string),
	}
	// systemctl show --no-page ssh
	sdCmd := exec.Command("systemctl", "show", "--no-page", unit)
	var stdo bytes.Buffer
	var stde bytes.Buffer
	sdCmd.Stdout = &stdo
	sdCmd.Stderr = &stde
	if err = sdCmd.Run(); err != nil {
		fmt.Println(stdo.String())
		fmt.Println(stde.String())
		return
	}

	r := bufio.NewScanner(&stdo)

	for r.Scan() {
		kv := strings.Split(r.Text(), "=")
		if len(kv) > 2 {
			continue
		}
		s.data[kv[0]] = kv[1]
	}

	fmt.Println(stdo.String())
	return
}

func (s Service) String() string {
	return fmt.Sprintf("%s state: %s preset: %s ", s.name, s.ActiveState(), s.Preset())
}

func (s Service) State() string {
	return s.data["UnitFileState"]
}

func (s Service) Preset() string {
	return s.data["UnitFilePreset"]
}

func (s Service) ActiveState() string {
	return s.data["ActiveState"]
}

func (s Service) ActiveStateInt() int {
	switch s.data["ActiveState"] {
	case "active":
		return 1
	case "activating":
		return 0
	case "reloading":
		return 0
	case "maintenance":
		return 0
	case "deactivating":
		return -1
	case "inactive":
		return -2
	case "failed":
		return -3
	}
	return -7
}

/*
active
activating
reloading
maintenance
deactivating
inactive
failed
*/
