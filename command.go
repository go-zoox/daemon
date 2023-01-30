package daemon

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/go-zoox/logger"
)

func Run(cfg *Config, cmdPath string, args ...string) error {
	realCmd := &exec.Cmd{
		Path: cmdPath,
		Args: append([]string{cmdPath}, args...),
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
