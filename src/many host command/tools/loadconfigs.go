package tools

import (
	"bufio"
	"os"
	"strings"
)

// HostInfo represents the configuration for a single remote host.
// It includes connection details and authentication credentials.
type HostInfo struct {
	IP       string // Remote host IP address
	Port     string // Remote host SSH port
	User     string // SSH username
	Password string // SSH password or path to private key file (legacy support)
	KeyPath  string // Path to SSH private key file
	Group    string // Group name the host belongs to
}

// ReadHosts reads host configurations from the specified file.
// It supports two formats:
// 1. Legacy: "IP User Password"
// 2. New: INI-style groups and key-value pairs (e.g., "address=... user=...")
//
// filename: The path to the host list file.
// targetGroup: If specified, only hosts belonging to this group will be returned.
// Returns a slice of HostInfo and any error encountered.
func ReadHosts(filename string, targetGroup string) ([]HostInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var hosts []HostInfo
	currentGroup := "default"
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 处理组名 [group]
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentGroup = line[1 : len(line)-1]
			continue
		}

		// 如果指定了组，且当前组不匹配，则跳过
		if targetGroup != "" && currentGroup != targetGroup {
			continue
		}

		// 解析行内容
		if strings.Contains(line, "=") {
			// 新格式: address=127.0.0.1 user=root password="" ssh_key=""
			host := HostInfo{Group: currentGroup}
			fields := strings.Fields(line)
			for _, field := range fields {
				kv := strings.SplitN(field, "=", 2)
				if len(kv) != 2 {
					continue
				}
				key := kv[0]
				value := strings.Trim(kv[1], "\"") // 去除引号

				switch key {
				case "address":
					host.IP = value
				case "port":
					host.Port = value
				case "user":
					host.User = value
				case "password":
					host.Password = value
				case "ssh_key":
					host.KeyPath = value
				}
			}
			if host.IP != "" {
				if host.User == "" {
					host.User = "root"
				}
				if host.Port == "" {
					host.Port = "22"
				}
				hosts = append(hosts, host)
			}
		} else {
			// 旧格式: IP User Password
			fields := strings.Fields(line)
			if len(fields) < 3 {
				continue
			}
			ip := fields[0]
			port := "22"
			if strings.Contains(ip, ":") {
				parts := strings.Split(ip, ":")
				ip = parts[0]
				port = parts[1]
			}
			hosts = append(hosts, HostInfo{
				IP:       ip,
				Port:     port,
				User:     fields[1],
				Password: fields[2],
				Group:    currentGroup,
			})
		}
	}
	return hosts, scanner.Err()
}
