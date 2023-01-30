package daemon

// reference: https://zhuanlan.zhihu.com/p/146192035

import (
	"fmt"
	"os"
	"os/exec"
)

type BackgroundConfig struct {
	Cmd     string
	Args    []string
	Role    string
	LogFile string
}

func Background(cfg *BackgroundConfig) (*exec.Cmd, error) {
	role := os.Getenv(EnvName)
	if role != "" {
		// 子进程 => 退出
		return nil, nil
	}

	// 父进程
	cmd := &exec.Cmd{
		Path: cfg.Cmd,
		Args: append([]string{cfg.Cmd}, cfg.Args...),
		Env:  os.Environ(),
	}

	// 为子进程设置特殊的环境变量标识
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", EnvName, cfg.Role))

	if cfg.LogFile != "" {
		stdout, err := os.OpenFile(cfg.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return nil, fmt.Errorf("[pid: %d] failed to open log file: %v", os.Getpid(), err)
		}

		cmd.Stderr = stdout
		cmd.Stdout = stdout
	}

	// 一步启动子进程
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return cmd, nil
}
