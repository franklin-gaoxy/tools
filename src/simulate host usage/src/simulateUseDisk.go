package src

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// SimulateUseDisk simulates disk usage by writing files to a specified path
func SimulateUseDisk(cmd *cobra.Command, args []string) error {
	// Retrieve command-line flags
	sizeGB, _ := cmd.Flags().GetInt("size")
	path, _ := cmd.Flags().GetString("path")
	durationMin, _ := cmd.Flags().GetInt("time")

	// Convert size from GB to bytes
	sizeBytes := int64(sizeGB) * 1024 * 1024 * 1024

	// Ensure the directory exists
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Channel to signal stop
	stopChan := make(chan struct{})
	var wg sync.WaitGroup

	// Start writing files in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		var prevFile string
		for {
			select {
			case <-stopChan:
				return
			default:
				fileName := filepath.Join(path, fmt.Sprintf("temp_%d", time.Now().UnixNano()))
				if err := writeFile(fileName, sizeBytes); err != nil {
					fmt.Printf("Error writing file: %v\n", err)
					continue
				}
				if prevFile != "" {
					go deleteFile(prevFile) // Asynchronously delete the previous file
				}
				prevFile = fileName
			}
		}
	}()

	// Set a timer to stop the program if duration is specified
	if durationMin > 0 {
		time.AfterFunc(time.Duration(durationMin)*time.Minute, func() {
			close(stopChan)
		})
	}

	// Wait for the stop signal and clean up
	<-stopChan
	wg.Wait()
	cleanUp(path)

	return nil
}

// writeFile writes a file of the specified size with arbitrary data
func writeFile(fileName string, sizeBytes int64) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	data := make([]byte, 1024) // 1KB buffer of arbitrary data
	for i := range data {
		data[i] = byte(i % 256)
	}

	written := int64(0)
	for written < sizeBytes {
		n, err := file.Write(data)
		if err != nil {
			return err
		}
		written += int64(n)
	}
	return nil
}

// deleteFile removes a specified file
func deleteFile(fileName string) {
	if err := os.Remove(fileName); err != nil {
		fmt.Printf("Error deleting file: %v\n", err)
	}
}

// cleanUp removes all files in the specified directory
func cleanUp(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}
	for _, file := range files {
		if !file.IsDir() {
			deleteFile(filepath.Join(path, file.Name()))
		}
	}
}
