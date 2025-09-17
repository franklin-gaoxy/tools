package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"many/printline"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

type HostInfo struct {
	IP       string
	User     string
	Password string
}

// 读取主机文件
func readHosts(filename string) ([]HostInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var hosts []HostInfo
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		hosts = append(hosts, HostInfo{
			IP:       fields[0],
			User:     fields[1],
			Password: fields[2],
		})
	}
	return hosts, scanner.Err()
}

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

func main() {
	var file string
	var batchSize int

	rootCmd := &cobra.Command{
		Use:   "ssh-tool [command]",
		Short: "SSH execution tool",
		Args:  cobra.MinimumNArgs(1), // 至少需要一个命令
		Run: func(cmd *cobra.Command, args []string) {
			hosts, err := readHosts(file)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			runParallel(hosts, strings.Join(args, " "))
		},
	}

	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "./nodelist", "hosts file")

	parallelCmd := &cobra.Command{
		Use:   "parallel [command]",
		Short: "Run in parallel mode",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hosts, err := readHosts(file)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			runParallel(hosts, strings.Join(args, " "))
		},
	}

	serialCmd := &cobra.Command{
		Use:   "serial [command]",
		Short: "Run in serial mode",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hosts, err := readHosts(file)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			runSerial(hosts, strings.Join(args, " "))
		},
	}

	batchCmd := &cobra.Command{
		Use:   "batch [command]",
		Short: "Run in batch mode",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hosts, err := readHosts(file)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			runBatch(hosts, strings.Join(args, " "), batchSize)
		},
	}
	batchCmd.Flags().IntVarP(&batchSize, "number", "n", 5, "batch size")

	rootCmd.AddCommand(parallelCmd, serialCmd, batchCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
