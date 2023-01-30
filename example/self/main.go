package main

import (
	"log"

	"github.com/go-zoox/daemon"
)

func main() {
	err := daemon.New(&daemon.Config{
		LogFile: "/tmp/gd.log",
		PidFile: "/tmp/gd.pid",
		Cmd:     "/usr/bin/top",
	}, func(cfg *daemon.Config) error {
		return daemon.Run(cfg, cfg.Cmd)
	})

	if err != nil {
		log.Fatal("daemonrized err:", err)
	}
}
