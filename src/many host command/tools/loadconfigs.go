package tools

import (
	"bufio"
	"os"
	"strings"
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
