package tools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/sftp"
	"many/tools/printline"
)

// RunScp copies files/directories to multiple hosts in parallel.
func RunScp(hosts []HostInfo, src, dest string) {
	var wg sync.WaitGroup
	for _, h := range hosts {
		wg.Add(1)
		go func(host HostInfo) {
			defer wg.Done()
			err := copyToHost(host, src, dest)
			header := fmt.Sprintf("[%s]", host.IP)
			if err != nil {
				printline.ExecuteCenter(header, "=", "y", "n")
				fmt.Printf("SCP failed: %v\n", err)
			} else {
				printline.ExecuteCenter(header, "=", "y", "n")
				fmt.Printf("SCP success: %s -> %s\n", src, dest)
			}
		}(h)
	}
	wg.Wait()
}

func copyToHost(host HostInfo, src, dest string) error {
	client, err := ConnectSSH(host)
	if err != nil {
		return err
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("sftp client error: %v", err)
	}
	defer sftpClient.Close()

	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat src error: %v", err)
	}

	if info.IsDir() {
		return copyDir(sftpClient, src, dest)
	}
	return copyFile(sftpClient, src, dest)
}

func copyFile(client *sftp.Client, localPath, remotePath string) error {
	srcFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Check if remotePath is a directory
	remoteInfo, err := client.Stat(remotePath)
	finalRemotePath := remotePath
	if err == nil && remoteInfo.IsDir() {
		finalRemotePath = filepath.Join(remotePath, filepath.Base(localPath))
	}

	dstFile, err := client.Create(finalRemotePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Preserve permissions
	if info, err := os.Stat(localPath); err == nil {
		client.Chmod(finalRemotePath, info.Mode())
	}
	return nil
}

func copyDir(client *sftp.Client, localPath, remotePath string) error {
	localName := filepath.Base(localPath)
	remoteInfo, err := client.Stat(remotePath)
	targetDir := remotePath

	// If remote path exists and is a directory, copy INTO it
	if err == nil && remoteInfo.IsDir() {
		targetDir = filepath.Join(remotePath, localName)
	}

	// Create target dir
	if err := client.MkdirAll(targetDir); err != nil {
		return err
	}

	return filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(localPath, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}

		remoteDest := filepath.Join(targetDir, relPath)

		if info.IsDir() {
			return client.MkdirAll(remoteDest)
		}

		// Direct file copy without checking if remoteDest is dir (it shouldn't be)
		return copyFileDirect(client, path, remoteDest)
	})
}

// copyFileDirect copies to a specific path, assuming it's the filename
func copyFileDirect(client *sftp.Client, localPath, remotePath string) error {
	srcFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := client.Create(remotePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	if info, err := os.Stat(localPath); err == nil {
		client.Chmod(remotePath, info.Mode())
	}
	return nil
}
