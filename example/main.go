package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/go-zoox/daemon"
)

func main() {
	err := daemon.Daemon(&daemon.Config{
		LogFile: "/tmp/gd.log",
		PidFile: "/tmp/gd.pid",
	}, func(cfg *daemon.Config) *exec.Cmd {
		cmd := "/usr/bin/top"
		return &exec.Cmd{
			Path: cmd,
			Args: []string{cmd},
			Env:  os.Environ(),
		}
	})

	if err != nil {
		log.Fatal("daemonrized err:", err)
	}
}
