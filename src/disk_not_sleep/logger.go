package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"k8s.io/klog/v2"
)

// newLogger 初始化 klog，并把日志输出到 stdout 以及可选的日志文件。
//
// 约定：每次启动都会创建一个新的日志文件，文件名为：
//
//	<log_file_prefix>_<启动时间>.log
//
// 启动时间格式使用精确到毫秒的时间戳，避免同一秒内多次启动发生冲突。
func newLogger(cfg DiskNotSleepConfig) (func(), string, error) {
	klog.InitFlags(nil)
	_ = flag.Set("logtostderr", "false")
	_ = flag.Set("alsologtostderr", "false")
	_ = flag.Set("v", verbosityFromLevel(cfg.LogLevel))

	writer := io.Writer(os.Stdout)
	closeFn := func() { klog.Flush() }
	logFile := ""

	logDir := strings.TrimSpace(cfg.LogFilePath)
	if logDir == "" {
		klog.SetOutput(writer)
		return closeFn, logFile, nil
	}

	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return func() {}, "", fmt.Errorf("create log_file_path: %w", err)
	}

	prefix := strings.TrimSpace(cfg.LogFilePrefix)
	if prefix == "" {
		prefix = "disk_not_sleep"
	}

	name := fmt.Sprintf("%s_%s.log", prefix, time.Now().Format("20060102_150405_000"))
	fullPath := filepath.Join(logDir, name)

	f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return func() {}, "", fmt.Errorf("open log file: %w", err)
	}

	writer = io.MultiWriter(os.Stdout, f)
	klog.SetOutput(writer)
	logFile = fullPath

	closeFn = func() {
		klog.Flush()
		_ = f.Close()
	}
	return closeFn, logFile, nil
}

func verbosityFromLevel(level string) string {
	level = strings.ToLower(strings.TrimSpace(level))
	switch level {
	case "debug":
		return "4"
	case "info":
		return "2"
	case "warn", "warning":
		return "1"
	case "error":
		return "0"
	default:
		return "2"
	}
}
