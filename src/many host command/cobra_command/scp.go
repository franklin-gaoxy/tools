package cobra_command

import (
	"fmt"
	"io"
	"many/tools"
	"many/tools/printline"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/sftp"
	"github.com/spf13/cobra"
)

// NewScpCommand creates a command to copy files to target hosts.
func NewScpCommand() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "scp [src] [dest]",
		Short: "Copy files to remote hosts",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			hosts, err := GetHosts()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			src := args[0]
			dest := args[1]
			runScp(hosts, src, dest, limit)
		},
	}
	cmd.Flags().StringVarP(&Group, "group", "g", "", "target group (default select all)")
	cmd.Flags().IntVarP(&limit, "number", "n", 0, "concurrent limit (0 for unlimited)")
	return cmd
}

// runScp copies files/directories to multiple hosts in parallel.
func runScp(hosts []tools.HostInfo, src, dest string, limit int) {
	var wg sync.WaitGroup

	// 如果 limit > 0，创建一个 worker pool channel 来限制并发
	var sem chan struct{}
	if limit > 0 {
		sem = make(chan struct{}, limit)
	}

	for _, h := range hosts {
		// 如果启用了并发限制，获取一个 token
		if limit > 0 {
			sem <- struct{}{}
		}

		wg.Add(1)
		go func(host tools.HostInfo) {
			defer wg.Done()
			// 任务完成后释放 token
			if limit > 0 {
				defer func() { <-sem }()
			}

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

func copyToHost(host tools.HostInfo, src, dest string) error {
	client, err := tools.ConnectSSH(host)
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
