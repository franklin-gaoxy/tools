package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

// defaultConfigPath：当未指定 -f 时，仅在当前工作目录查找的默认配置文件名。
const defaultConfigPath = "disk_not_sleep.yaml"

func main() {
	ctx, stop := notifyContext(context.Background())
	defer stop()

	cmd := newRootCmd()
	cmd.SetContext(ctx)
	if err := cmd.Execute(); err != nil {
		showFatalError(err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var configPath string
	var enableTray bool

	cmd := &cobra.Command{
		Use: "disk_not_sleep",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, loadedFrom, err := resolveConfig(configPath)
			if err != nil {
				return err
			}

			closeLogger, logFile, err := newLogger(cfg)
			if err != nil {
				return err
			}
			defer closeLogger()

			if loadedFrom != "" {
				klog.Infof("loaded config: %s", loadedFrom)
			} else {
				klog.Infof("loaded config: built-in defaults")
			}
			if logFile != "" {
				klog.Infof("log file: %s", logFile)
			} else {
				klog.Infof("log file: stdout only")
			}
			klog.Infof("tmp write target: dir=%s name=%s", cfg.TmpFilePath, cfg.TmpFileName)
			klog.Infof("time interval: %s", cfg.TimeInterval)

			runCtx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			if enableTray {
				errCh := make(chan error, 1)
				go func() {
					errCh <- runDiskNotSleep(runCtx, cfg)
				}()
				runWithTray(runCtx, cancel, cfg, loadedFrom, logFile)
				err := <-errCh
				if errors.Is(err, context.Canceled) {
					return nil
				}
				return err
			}

			if err := runDiskNotSleep(runCtx, cfg); err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "f", "", "config file path")
	cmd.Flags().BoolVar(&enableTray, "tray", true, "run minimized in system tray")

	cmd.SetHelpTemplate(fmt.Sprintf("%s\n\nFlags:\n  -f, --config string   config file path\n\nConfig keys:\n  tmp_file_path: directory to create/write keep-alive file\n  tmp_file_name: keep-alive file name\n  time_interval: duration like 180s, 5m\n  log_level: debug|info|warn|error\n  log_file_path: directory to write log file\n  log_file_prefix: log file name prefix\n", cmd.HelpTemplate()))

	return cmd
}

func resolveConfig(flagPath string) (DiskNotSleepConfig, string, error) {
	// 配置优先级：
	// 1) -f 指定路径
	// 2) 可执行文件同级目录下的 disk_not_sleep.yaml
	// 3) 当前工作目录下的 disk_not_sleep.yaml
	// 3) 内置默认配置
	if flagPath != "" {
		cfg, err := loadConfigFromFile(flagPath)
		if err != nil {
			return DiskNotSleepConfig{}, "", err
		}
		if err := cfg.Validate(); err != nil {
			return DiskNotSleepConfig{}, "", err
		}
		return cfg, flagPath, nil
	}

	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	if exeDir != "" {
		p := filepath.Join(exeDir, defaultConfigPath)
		if _, err := os.Stat(p); err == nil {
			cfg, err := loadConfigFromFile(p)
			if err != nil {
				return DiskNotSleepConfig{}, "", err
			}
			if err := cfg.Validate(); err != nil {
				return DiskNotSleepConfig{}, "", err
			}
			return cfg, p, nil
		}
	}

	if _, err := os.Stat(defaultConfigPath); err == nil {
		cfg, err := loadConfigFromFile(defaultConfigPath)
		if err != nil {
			return DiskNotSleepConfig{}, "", err
		}
		if err := cfg.Validate(); err != nil {
			return DiskNotSleepConfig{}, "", err
		}
		return cfg, defaultConfigPath, nil
	}

	cfg := defaultConfig()
	if err := cfg.Validate(); err != nil {
		return DiskNotSleepConfig{}, "", err
	}
	return cfg, "", nil
}
