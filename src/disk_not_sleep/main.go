package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	var (
		path     string
		filename string
		interval int
	)

	var rootCmd = &cobra.Command{
		Use:   "disk_not_sleep",
		Short: "A tool to prevent disk sleep by writing temp files periodically",
		Run: func(cmd *cobra.Command, args []string) {
			// Validate path
			if path == "" {
				// Default to current directory if not specified
				path = "."
			}

			// Ensure directory exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Printf("Error: Directory '%s' does not exist.\n", path)
				os.Exit(1)
			}

			fmt.Printf("Starting disk_not_sleep...\n")
			fmt.Printf("Target Path: %s\n", path)
			fmt.Printf("Filename: %s\n", filename)
			fmt.Printf("Interval: %d seconds\n", interval)

			ticker := time.NewTicker(time.Duration(interval) * time.Second)
			defer ticker.Stop()

			// Run once immediately
			keepAlive(path, filename)

			for range ticker.C {
				keepAlive(path, filename)
			}
		},
	}

	rootCmd.Flags().StringVarP(&path, "path", "p", ".", "Target directory path")
	rootCmd.Flags().StringVarP(&filename, "filename", "f", "tmp_file", "Temporary file name")
	rootCmd.Flags().IntVarP(&interval, "time", "t", 60, "Interval time in seconds")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func keepAlive(dir, name string) {
	fullPath := filepath.Join(dir, name)

	// Create/Overwrite file
	f, err := os.Create(fullPath)
	if err != nil {
		fmt.Printf("[%s] Error creating file: %v\n", time.Now().Format(time.RFC3339), err)
		return
	}

	// Write "0"
	_, err = f.WriteString("0")
	if err != nil {
		f.Close() // Close before printing error
		fmt.Printf("[%s] Error writing file: %v\n", time.Now().Format(time.RFC3339), err)
		return
	}
	f.Close() // Close immediately

	// Delete file
	err = os.Remove(fullPath)
	if err != nil {
		fmt.Printf("[%s] Error deleting file: %v\n", time.Now().Format(time.RFC3339), err)
		return
	}

	fmt.Printf("[%s] Keep-alive cycle completed: Wrote & Deleted %s\n", time.Now().Format(time.RFC3339), fullPath)
}
