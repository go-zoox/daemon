package model

import (
	"sync"

	"github.com/go-zoox/daemon"
)

type MasterConfig struct {
	daemon.Config
	MaxSalves int
}

func Master(cfg *MasterConfig) error {
	return daemon.Daemon(&daemon.Config{
		LogFile: cfg.LogFile,
		PidFile: cfg.PidFile,
	}, func(c *daemon.Config) error {
		wg := sync.WaitGroup{}
		pids := map[int]bool{}

		i := 0
		for {
			if i >= cfg.MaxSalves {
				break
			}
			i++

			cmd, err := Slave(&SlaceConfig{})
			if err != nil {
				return err
			}

			pids[cmd.Process.Pid] = true
		}

		wg.Wait()
		return nil
	})
}
