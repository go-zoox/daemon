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
	cmd = exec.Command(cfg.Cmd, cfg.Args...)

	cmd.Env = os.Environ()

	return cmd, nil
}
