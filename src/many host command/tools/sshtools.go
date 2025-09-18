package tools

import (
	"fmt"
	"io/ioutil"
	"log"
	"many/printline"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// 执行命令
func runCommand(host HostInfo, cmd string) (string, error) {
	var auth ssh.AuthMethod
	if _, err := os.Stat(host.Password); err == nil {
		// 当作私钥文件
		key, err := ioutil.ReadFile(host.Password)
		if err != nil {
			return "", err
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return "", err
		}
		auth = ssh.PublicKeys(signer)
	} else {
		// 当作密码
		auth = ssh.Password(host.Password)
	}

	config := &ssh.ClientConfig{
		User:            host.User,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host.IP+":22", config)
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

// 并行执行
func runParallel(hosts []HostInfo, cmd string) {
	var wg sync.WaitGroup
	for _, h := range hosts {
		wg.Add(1)
		go func(host HostInfo) {
			defer wg.Done()
			out, err := runCommand(host, cmd)
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

// 串行执行
func runSerial(hosts []HostInfo, cmd string) {
	for _, h := range hosts {
		out, err := runCommand(h, cmd)
		header := fmt.Sprintf("\n[%s]\n", h.IP)
		if err != nil {
			fmt.Printf("%s%s", header, out)
			fmt.Printf("[%s] Command failed: %v\n", h.IP, err)
		} else {
			fmt.Printf("%s%s", header, out)
		}
	}
}

// 批次执行
func runBatch(hosts []HostInfo, cmd string, batchSize int) {
	for i := 0; i < len(hosts); i += batchSize {
		end := i + batchSize
		if end > len(hosts) {
			end = len(hosts)
		}
		batch := hosts[i:end]
		runParallel(batch, cmd)
	}
}

// TestConnection 尝试连接一批主机，打印连接是否成功
func TestConnection(hosts []HostInfo) {
	for _, host := range hosts {
		addr := fmt.Sprintf("%s:22", host.IP)

		config := &ssh.ClientConfig{
			User: host.User,
			Auth: []ssh.AuthMethod{
				ssh.Password(host.Password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         5 * time.Second,
		}

		client, err := ssh.Dial("tcp", addr, config)
		if err != nil {
			log.Printf("[FAIL] %s@%s → %v\n", host.User, host.IP, err)
			continue
		}
		defer client.Close()

		log.Printf("[OK] %s@%s 连接成功\n", host.User, host.IP)
	}
}
