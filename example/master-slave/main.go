package main

import (
	"log"

	"github.com/go-zoox/daemon"
	"github.com/go-zoox/daemon/model"
)

func main() {
	if err := model.Run(&model.Config{
		Config: daemon.Config{
			PidFile: "/tmp/master-slave.pid",
			LogFile: "/tmp/master-slave.log",
		},
		PidLogFile: "/tmp/master-slave.pid.log",
		MaxSalves:  4,
		//
		// Cmd:       "/usr/local/bin/docker",
		//
		// Cmd:  "/tmp/mm.sh",
		// Args: []string{"stats"},
		//
		Cmd: "top",
	}); err != nil {
		log.Fatal(err)
	}
}
