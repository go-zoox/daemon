package daemon

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/go-zoox/fs"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/logger/components/transport"
	"github.com/go-zoox/logger/transport/file"
)

type Config struct {
	Cmd     string
	Args    []string
	Role    string
	LogFile string
	PidFile string
}

func Daemon(cfg *Config) error {
	logger.SetTransports(map[string]transport.Transport{
		"file": file.New(&file.Config{
			Level: "info",
			File:  cfg.LogFile,
		}),
	})

	cmd, err := Background(&BackgroundConfig{
		Role:    RoleMaster,
		LogFile: cfg.LogFile,
		Cmd:     os.Args[0],
		Args:    os.Args[1:], // 注意,此处是包含程序名的
	})
	if err != nil {
		return fmt.Errorf("failed to start daemon master: %v", err)
	}

	if cmd != nil {
		// parent => exit
		return nil
	}

	if cfg.PidFile != "" {
		// @TODO
		if fs.IsExist(cfg.PidFile) {
			if err := fs.RemoveFile(cfg.PidFile); err != nil {
				return fmt.Errorf("[pid: %d] failed to remove pid file(%s): %v", os.Getpid(), cfg.PidFile, err)
			}
		}

		if err := fs.WriteFile(cfg.PidFile, []byte(fmt.Sprintf("%d", os.Getpid()))); err != nil {
			return fmt.Errorf("[pid: %d] failed to write pid to file(%s): %v", os.Getpid(), cfg.PidFile, err)
		}
	}

	logger.Infof("[daemon: %d] start ...", os.Getpid())

	realCmd := &exec.Cmd{
		Path: cfg.Cmd,
		Args: append([]string{cfg.Cmd}, cfg.Args...),
		Env:  os.Environ(),
	}

	if cfg.LogFile != "" {
		stdout, err := os.OpenFile(cfg.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return fmt.Errorf("[pid: %d] failed to open log file: %v", os.Getpid(), err)
		}

		realCmd.Stderr = stdout
		realCmd.Stdout = stdout
	}

	if err := realCmd.Start(); err != nil {
		return fmt.Errorf("[daemon: %d] failed to start daemon: %v", os.Getpid(), err)
	}

	if err := realCmd.Wait(); err != nil {
		return fmt.Errorf("[daemon: %d] failed to wait daemon: %v", os.Getpid(), err)
	}

	logger.Infof("[daemon: %d] exit", os.Getpid())
	return nil
}
