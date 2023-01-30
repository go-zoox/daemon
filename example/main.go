package main

import (
	"log"

	"github.com/go-zoox/daemon"
)

func main() {
	err := daemon.Daemon(&daemon.Config{
		LogFile: "/tmp/gd.log",
		PidFile: "/tmp/gd.pid",
	}, func(cfg *daemon.Config) error {
		return daemon.RunCommand(cfg, "/usr/bin/top")
	})

	if err != nil {
		log.Fatal("daemonrized err:", err)
	}
}
