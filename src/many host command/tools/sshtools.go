package tools

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

// ConnectSSH establishes an SSH connection to the host.
func ConnectSSH(host HostInfo) (*ssh.Client, error) {
	var authMethods []ssh.AuthMethod

	// 1. 优先尝试使用 SSH Key (指定了 ssh_key 字段)
	if host.KeyPath != "" {
		key, err := ioutil.ReadFile(host.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("ssh key file read error: %v", err)
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("ssh key parse error: %v", err)
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
	if !isLegacyKey && host.Password != "" {
		authMethods = append(authMethods, ssh.Password(host.Password))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no authentication methods provided (password or ssh_key)")
	}

	config := &ssh.ClientConfig{
		User:            host.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%s", host.IP, host.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("ssh connect error: %v", err)
	}
	return client, nil
}

// RunCommand executes a single command on a remote host via SSH.
// It handles authentication using either SSH keys or passwords.
// If both are provided, it prioritizes the SSH key but will fall back to the password.
//
// host: The target host configuration.
// cmd: The command string to execute.
// Returns the combined stdout/stderr output and any error.
func RunCommand(host HostInfo, cmd string) (string, error) {
	client, err := ConnectSSH(host)
	if err != nil {
		return "", err
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
