package tools

import (
	"fmt"
	"io/ioutil"
	"many/tools/printline"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
)

// RunCommand executes a single command on a remote host via SSH.
// It handles authentication using either SSH keys or passwords.
// If both are provided, it prioritizes the SSH key but will fall back to the password.
//
// host: The target host configuration.
// cmd: The command string to execute.
// Returns the combined stdout/stderr output and any error.
func RunCommand(host HostInfo, cmd string) (string, error) {
	var authMethods []ssh.AuthMethod

	// 1. 优先尝试使用 SSH Key (指定了 ssh_key 字段)
	if host.KeyPath != "" {
		key, err := ioutil.ReadFile(host.KeyPath)
		if err != nil {
			return "", fmt.Errorf("ssh key file read error: %v", err)
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return "", fmt.Errorf("ssh key parse error: %v", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// 2. 处理 Password 字段
	// 兼容旧逻辑：如果 ssh_key 为空，且 password 字段指向一个文件，则尝试作为私钥加载
	isLegacyKey := false
	if host.KeyPath == "" && host.Password != "" {
		if info, err := os.Stat(host.Password); err == nil && !info.IsDir() {
			key, err := ioutil.ReadFile(host.Password)
			if err == nil {
				signer, err := ssh.ParsePrivateKey(key)
				if err == nil {
					authMethods = append(authMethods, ssh.PublicKeys(signer))
					isLegacyKey = true
				}
			}
		}
	}

	// 如果不是作为旧版私钥文件处理，则作为普通密码添加
	// 这样如果同时指定了 ssh_key 和 password，两者都会被加入（key 在前）
	if !isLegacyKey && host.Password != "" {
		authMethods = append(authMethods, ssh.Password(host.Password))
	}

	if len(authMethods) == 0 {
		return "", fmt.Errorf("no authentication methods provided (password or ssh_key)")
	}

	config := &ssh.ClientConfig{
		User:            host.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%s", host.IP, host.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return "", fmt.Errorf("ssh connect error: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("new session error: %v", err)
	}
	defer session.Close()

	// stdout + stderr
	output, err := session.CombinedOutput(cmd)
	return string(output), err
}

// RunParallel executes a command on multiple hosts concurrently.
// It uses a goroutine for each host and waits for all to complete.
func RunParallel(hosts []HostInfo, cmd string) {
	var wg sync.WaitGroup
	for _, h := range hosts {
		wg.Add(1)
		go func(host HostInfo) {
			defer wg.Done()
			out, err := RunCommand(host, cmd)
			header := fmt.Sprintf("[%s]", host.IP)
			if err != nil {
				// 即使报错，也打印 stderr+stdout
				printline.ExecuteCenter(header, "=", "y", "n")
				fmt.Println(out)
				fmt.Printf("Command failed: %v\n", err)
			} else {
				// fmt.Printf("%s%s", header, out)
				printline.ExecuteCenter(header, "=", "y", "n")
				// printline.ExecutePrintLine("=", "y")
				fmt.Println(out)
			}
		}(h)
	}
	wg.Wait()
}
