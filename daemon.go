package daemon

import (
	"fmt"
	"os"

	"github.com/go-zoox/fs"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/logger/components/transport"
	"github.com/go-zoox/logger/transport/file"
)

// Config is the Daemon Config.
type Config struct {
	LogFile string
	PidFile string
	//
	Cmd  string
	Args []string
}

// New daemonizes a command.
func New(cfg *Config, onRun func(cfg *Config) error) error {
	logger.SetTransports(map[string]transport.Transport{
		"file": file.New(&file.Config{
			Level: "info",
			File:  cfg.LogFile,
		}),
	})

	cmd, err := Background(&BackgroundConfig{
		LogFile: cfg.LogFile,
		Cmd:     cfg.Cmd,
		Args:    cfg.Args,
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

	// root
	logger.Infof("[daemon: %d][roor] start ...", os.Getpid())
	err = onRun(cfg)
	if err != nil {
		// return err
		logger.Infof("[daemon: %d][root] exit ...", os.Getpid())
	}

	return err

}

// Daemonrize daemonizes a command.
func Daemonrize(cfg *Config, onRun func(cfg *Config) error) error {
	return New(cfg, onRun)
}
