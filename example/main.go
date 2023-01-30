package main

import (
	"log"

	"github.com/go-zoox/daemon"
)

func main() {
	err := daemon.New(&daemon.Config{
		LogFile: "/tmp/gd.log",
		PidFile: "/tmp/gd.pid",
	}, func(cfg *daemon.Config) error {
		return daemon.Run(cfg, "/usr/bin/top")
	})

	if err != nil {
		log.Fatal("daemonrized err:", err)
	}
}
