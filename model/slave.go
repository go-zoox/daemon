package model

import "os/exec"

type SlaceConfig struct {
	Cmd  string
	Args []string
	ID   int
}

func Slave(cfg *SlaceConfig) (cmd *exec.Cmd, err error) {
	return
}
