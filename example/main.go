package main

import (
	"log"

	"github.com/go-zoox/daemon"
)

func main() {
	if err := daemon.Daemon(&daemon.Config{
		LogFile: "/tmp/gd.log",
		PidFile: "/tmp/gd.pid",
		Role:    daemon.RoleMaster,
		Cmd:     "/usr/bin/top",
	}); err != nil {
		log.Fatal("daemonrized err:", err)
	}
}
