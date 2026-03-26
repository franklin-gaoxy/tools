package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type DiskNotSleepConfig struct {
	// TmpFilePath：要在哪个目录下创建临时文件；不存在则创建。
	TmpFilePath string `yaml:"tmp_file_path"`
	// TmpFileName：临时文件名；最终写入路径为 filepath.Join(TmpFilePath, TmpFileName)。
	TmpFileName string `yaml:"tmp_file_name"`
	// TimeInterval：间隔时间，每隔多久写入一次（例如："180s"、"5m"）。
	TimeInterval string `yaml:"time_interval"`
	// LogLevel：日志级别（debug|info|warn|error）。
	LogLevel string `yaml:"log_level"`
	// LogFilePath：日志文件输出目录；为空则只输出到 stdout。
	LogFilePath string `yaml:"log_file_path"`
	// LogFilePrefix：日志文件名前缀；日志文件名为 <prefix>_<启动时间>.log。
	LogFilePrefix string `yaml:"log_file_prefix"`
}

type configWrapper struct {
	DiskNotSleep *DiskNotSleepConfig `yaml:"disk_not_sleep"`
}

func defaultConfig() DiskNotSleepConfig {
	return DiskNotSleepConfig{
		TmpFilePath:   "D:\\system disk",
		TmpFileName:   "disk_not_sleep",
		TimeInterval:  "180s",
		LogLevel:      "info",
		LogFilePath:   "D:\\system disk",
		LogFilePrefix: "disk_not_sleep",
	}
}

func loadConfigFromFile(path string) (DiskNotSleepConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return DiskNotSleepConfig{}, fmt.Errorf("read config file: %w", err)
	}

	var wrapped configWrapper
	if err := yaml.Unmarshal(b, &wrapped); err != nil {
		return DiskNotSleepConfig{}, fmt.Errorf("parse config yaml: %w", err)
	}
	if wrapped.DiskNotSleep != nil {
		return *wrapped.DiskNotSleep, nil
	}

	var direct DiskNotSleepConfig
	if err := yaml.Unmarshal(b, &direct); err == nil {
		if !isZeroConfig(direct) {
			return direct, nil
		}
	}

	return DiskNotSleepConfig{}, errors.New("config yaml missing 'disk_not_sleep' section")
}

func isZeroConfig(cfg DiskNotSleepConfig) bool {
	return cfg.TmpFilePath == "" && cfg.TmpFileName == "" && cfg.TimeInterval == "" && cfg.LogLevel == "" && cfg.LogFilePath == "" && cfg.LogFilePrefix == ""
}

func (c DiskNotSleepConfig) Validate() error {
	if strings.TrimSpace(c.TmpFilePath) == "" {
		return errors.New("tmp_file_path is required")
	}
	if strings.TrimSpace(c.TmpFileName) == "" {
		return errors.New("tmp_file_name is required")
	}
	if strings.TrimSpace(c.TimeInterval) == "" {
		return errors.New("time_interval is required")
	}
	if _, err := c.IntervalDuration(); err != nil {
		return err
	}
	return nil
}

func (c DiskNotSleepConfig) IntervalDuration() (time.Duration, error) {
	s := strings.TrimSpace(c.TimeInterval)
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("invalid time_interval '%s': %w", s, err)
	}
	if d <= 0 {
		return 0, fmt.Errorf("time_interval must be > 0, got %s", d)
	}
	return d, nil
}
