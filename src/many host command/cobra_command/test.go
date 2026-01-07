package cobra_command

import (
	"fmt"
	"io/ioutil"
	"many/tools"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"k8s.io/klog/v2"
)

// NewTestCommand creates a command to test SSH connectivity to target hosts.
func NewTestCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Run test connect",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// Test command 也应该支持 group 吗？用户没明确说，但通常工具行为一致比较好
			// 不过为了严格遵守用户需求“给batch、serial、parallel三个都增加子参数”，暂时不给test加
			// 但是 ReadHosts 需要两个参数了，这里如果不传 Group 变量，就是传空字符串，表示全部
			// 现在使用 GetHosts()，它会使用 Group 变量，这意味着 Test 也支持 Group 了，这很好
			hosts, err := GetHosts()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			testConnection(hosts)
		},
	}
}

// testConnection attempts to establish an SSH connection to each host.
// It logs the success or failure of each connection attempt.
func testConnection(hosts []tools.HostInfo) {
	for _, host := range hosts {
		addr := fmt.Sprintf("%s:%s", host.IP, host.Port)
		klog.V(2).Infof("Attempting to connect to %s...", addr)

		var authMethods []ssh.AuthMethod

		// 1. 优先尝试使用 SSH Key (指定了 ssh_key 字段)
		if host.KeyPath != "" {
			if _, err := os.Stat(host.KeyPath); err == nil {
				key, err := ioutil.ReadFile(host.KeyPath)
				if err == nil {
					signer, err := ssh.ParsePrivateKey(key)
					if err == nil {
						authMethods = append(authMethods, ssh.PublicKeys(signer))
					}
				}
			}
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
			klog.Infof("[FAIL] %s@%s → no valid authentication methods provided\n", host.User, host.IP)
			continue
		}

		config := &ssh.ClientConfig{
			User:            host.User,
			Auth:            authMethods,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         5 * time.Second,
		}

		client, err := ssh.Dial("tcp", addr, config)
		if err != nil {
			klog.Infof("[FAIL] %s@%s → %v\n", host.User, host.IP, err)
			continue
		}
		defer client.Close()

		klog.Infof("[OK] %s@%s 连接成功\n", host.User, host.IP)
	}
}
