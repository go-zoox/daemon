package daemon

// reference: https://zhuanlan.zhihu.com/p/146192035

import (
	"fmt"
	"os"
	"os/exec"
)

// BackgroundConfig ...
type BackgroundConfig struct {
	LogFile string
}

// Background runs cmd in background.
func Background(cfg *BackgroundConfig) (*exec.Cmd, error) {
	childMark := os.Getenv(EnvName)
	if childMark != "" {
		// 子进程 => 退出
		return nil, nil
	}

	// 父进程
	cmd := &exec.Cmd{
		Path: os.Args[0],
		Args: os.Args, // 注意,此处是包含程序名的
		Env:  os.Environ(),
	}

	// 为子进程设置特殊的环境变量标识
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", EnvName, "IS_GO_ZOOX_DAEMON_CHILD"))

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
