package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/getlantern/systray"
	"k8s.io/klog/v2"
)

func runWithTray(ctx context.Context, cancel context.CancelFunc, cfg DiskNotSleepConfig, loadedFrom string, logFile string) {
	systray.Run(func() {
		systray.SetTitle("disk_not_sleep")
		systray.SetTooltip("disk_not_sleep")

		if loadedFrom != "" {
			systray.AddMenuItem(fmt.Sprintf("配置: %s", loadedFrom), "").Disable()
		} else {
			systray.AddMenuItem("配置: built-in defaults", "").Disable()
		}
		if logFile != "" {
			systray.AddMenuItem(fmt.Sprintf("日志: %s", logFile), "").Disable()
		} else {
			systray.AddMenuItem("日志: stdout only", "").Disable()
		}
		systray.AddMenuItem(fmt.Sprintf("写入目标: %s", filepath.Join(cfg.TmpFilePath, cfg.TmpFileName)), "").Disable()
		systray.AddMenuItem(fmt.Sprintf("间隔: %s", cfg.TimeInterval), "").Disable()
		systray.AddSeparator()

		mQuit := systray.AddMenuItem("退出", "退出程序")

		go func() {
			select {
			case <-mQuit.ClickedCh:
				klog.Infof("systray quit clicked")
				cancel()
				systray.Quit()
			case <-ctx.Done():
				systray.Quit()
			}
		}()
	}, func() {
		cancel()
	})
}
