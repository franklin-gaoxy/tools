package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"k8s.io/klog/v2"
)

// runDiskNotSleep 通过定时写入一个极小的文件来阻止磁盘休眠。
//
// 注意：这里选择“写入”而不是“写入后删除”，目的是尽量减少文件系统元数据抖动。
// 对于多数外置盘/机械盘场景，周期性写入即可触发设备保持活跃。
func runDiskNotSleep(ctx context.Context, cfg DiskNotSleepConfig) error {
	interval, err := cfg.IntervalDuration()
	if err != nil {
		return err
	}

	path := strings.TrimSpace(cfg.TmpFilePath)
	name := strings.TrimSpace(cfg.TmpFileName)
	fullPath := filepath.Join(path, name)

	klog.Infof("disk_not_sleep started: tmp_file=%s time_interval=%s", fullPath, interval)

	writeOnce := func() error {
		if err := os.MkdirAll(path, 0o755); err != nil {
			return fmt.Errorf("create tmp_file_path: %w", err)
		}
		f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			return fmt.Errorf("open tmp file: %w", err)
		}
		// write 1
		_, werr := f.WriteString("1")
		syncErr := f.Sync()
		cerr := f.Close()
		if werr != nil {
			return fmt.Errorf("write tmp file: %w", werr)
		}
		if syncErr != nil {
			return fmt.Errorf("sync tmp file: %w", syncErr)
		}
		if cerr != nil {
			return fmt.Errorf("close tmp file: %w", cerr)
		}
		klog.Infof("keep-alive write: %s", fullPath)
		return nil
	}

	if err := writeOnce(); err != nil {
		klog.Errorf("keep-alive write failed: %v", err)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			klog.Infof("disk_not_sleep stopped")
			return ctx.Err()
		case <-ticker.C:
			if err := writeOnce(); err != nil {
				klog.Errorf("keep-alive write failed: %v", err)
			}
		}
	}
}
