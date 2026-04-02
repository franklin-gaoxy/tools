package main

import (
	"os"
	"path/filepath"
)

func defaultLogDir() string {
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		if exeDir != "" {
			p := filepath.Join(exeDir, "logs")
			return p
		}
	}
	return filepath.Join(os.TempDir(), "disk_not_sleep_logs")
}
