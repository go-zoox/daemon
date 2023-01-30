package daemon

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/go-zoox/logger"
)

type Config struct {
	BackgroundConfig
}

func Daemon(cfg *Config) error {
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

	logger.Infof("[daemon: %d] start ...", os.Getegid())

	realCmd := &exec.Cmd{
		Path: cfg.Cmd,
		Args: append([]string{cfg.Cmd}, cfg.Args...),
		Env:  os.Environ(),
	}
	if err := realCmd.Start(); err != nil {
		return fmt.Errorf("[daemon: %d] failed to start daemon: %v", os.Getegid(), err)
	}

	if err := realCmd.Wait(); err != nil {
		return fmt.Errorf("[daemon: %d] failed to wait daemon: %v", os.Getegid(), err)
	}

	logger.Infof("[daemon: %d] exit", os.Getegid())
	return nil
}
