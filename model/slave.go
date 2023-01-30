package model

import (
	"os"
	"os/exec"
)

type SlaceConfig struct {
	Cmd  string
	Args []string
	ID   int
}

func Slave(cfg *SlaceConfig) (cmd *exec.Cmd, err error) {
	cmd = &exec.Cmd{
		Path: cfg.Cmd,
		Args: append([]string{cfg.Cmd}, cfg.Args...),
		Env:  os.Environ(),
	}

	return cmd, nil
}
