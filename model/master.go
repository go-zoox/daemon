package model

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-zoox/daemon"
	"github.com/go-zoox/fs"
)

type MasterConfig struct {
	daemon.Config
	MaxSalves  int
	PidLogFile string
	Cmd        string
	Args       []string
}

func Master(cfg *MasterConfig) error {
	return daemon.New(&daemon.Config{
		LogFile: cfg.LogFile,
		PidFile: cfg.PidFile,
	}, func(c *daemon.Config) error {
		var log *os.File
		var pidLog *os.File
		var err error
		cmds := map[int]*exec.Cmd{}
		stopped := false

		if cfg.LogFile != "" {
			log, err = os.OpenFile(cfg.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				return fmt.Errorf("[pid: %d] failed to open log file: %v", os.Getpid(), err)
			}
			defer log.Close()
		}

		if cfg.PidLogFile != "" {
			// remove old files
			if fs.IsExist(cfg.PidLogFile) {
				fs.RemoveFile(cfg.PidLogFile)
			}

			pidLog, err = os.OpenFile(cfg.PidLogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				return fmt.Errorf("[pid: %d] failed to open pid log file: %v", os.Getpid(), err)
			}
			defer pidLog.Close()
		} else {
			pidLog = log
		}

		pidLog.Write([]byte(fmt.Sprintf("[master: %d] started\n", os.Getpid())))

		// wait signal
		sigCh := make(chan os.Signal, 1)
		signal.Notify(
			sigCh,
			syscall.SIGTERM, // kill
			syscall.SIGHUP,
			syscall.SIGQUIT,
			syscall.SIGUSR1,
		)
		go func() {
			signal := <-sigCh
			stopped = true

			pidLog.Write([]byte(fmt.Sprintf("\n[master: %d] stopped with signal(%d)\n", os.Getpid(), signal)))
			for _, cmd := range cmds {
				if pidLog != nil {
					pidLog.Write([]byte(fmt.Sprintf("[slave: %d] kill by master\n", cmd.Process.Pid)))
				}

				go cmd.Process.Kill()
			}
		}()

		running := 0
		for {
			if stopped {
				break
			}

			if running >= cfg.MaxSalves {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			running += 1

			go func() {
				cmd, err := Slave(&SlaceConfig{
					Cmd:  cfg.Cmd,
					Args: cfg.Args,
				})
				if err != nil {
					return
				}

				cmd.Stderr = log
				cmd.Stdout = log
				if err = cmd.Start(); err != nil {
					return
				}

				pid := cmd.Process.Pid
				cmds[pid] = cmd
				pidLog.Write([]byte(fmt.Sprintf("[slave: %d] started\n", pid)))

				err = cmd.Wait()

				pidLog.Write([]byte(fmt.Sprintf("[slave: %d] stopped with error: %v\n", pid, err)))

				delete(cmds, pid)

				running -= 1
			}()
		}

		time.Sleep(100 * time.Millisecond)

		return nil
	})
}
