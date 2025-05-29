package src

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/mem"
	"github.com/spf13/cobra"
)

/*
simulate use memory
支持通过百分比或者指定大小 如10G 5G 或者 0.8 0.5
*/

// SimulateUseMemory 模拟内存使用情况
func SimulateUseMemory(cmd *cobra.Command, args []string) error {
	// 获取变量
	sizeStr, _ := cmd.Flags().GetString("size")
	duration, _ := cmd.Flags().GetInt("time")

	// Parse size into bytes
	targetBytes, err := parseSize(sizeStr)
	if err != nil {
		return fmt.Errorf("invalid size: %v", err)
	}

	// Allocate memory
	data, err := allocateMemory(targetBytes)
	if err != nil {
		return fmt.Errorf("failed to allocate memory: %v", err)
	}
	fmt.Printf("Allocated %d bytes of memory\n", targetBytes)

	// Handle duration if specified
	if duration > 0 {
		fmt.Printf("Running for %d minutes\n", duration)
		time.AfterFunc(time.Duration(duration)*time.Minute, func() {
			releaseMemory(&data)
			fmt.Println("Memory released, exiting")
			os.Exit(0)
		})
	}

	// Keep program running
	select {}
}

// parseSize 确定是百分比还是固定大小 返回需要创建的字节大小
func parseSize(sizeStr string) (uint64, error) {
	if isPercentage(sizeStr) {
		return parsePercentage(sizeStr)
	}
	return parseSpecificSize(sizeStr)
}

// isPercentage 检查输入是否为百分比
func isPercentage(sizeStr string) bool {
	return regexp.MustCompile(`^\d*\.?\d+$`).MatchString(sizeStr)
}

// parsePercentage 将总内存百分比转化为字节
func parsePercentage(sizeStr string) (uint64, error) {
	percent, err := strconv.ParseFloat(sizeStr, 64)
	// 范围验证
	if err != nil || percent <= 0 || percent > 1 {
		return 0, fmt.Errorf("percentage must be between 0 and 1")
	}
	vm, err := mem.VirtualMemory()
	if err != nil {
		return 0, fmt.Errorf("failed to get memory stats: %v", err)
	}
	// 当前使用加总数百分比
	currentUsage := vm.Used
	target := uint64(float64(vm.Total) * percent)
	// 如果现在主机运行使用的内存已经超过了指定的百分比 退出
	if target <= currentUsage {
		return 0, nil
	}
	// 因为运行期发现总会超出约200M左右的内存导致百分比不够精准 所以再减去200M
	return target - currentUsage - 209715200, nil
}

// parseSpecificSize 将指定大小的转换为字节
func parseSpecificSize(sizeStr string) (uint64, error) {
	re := regexp.MustCompile(`^(\d+)([GM])$`)
	matches := re.FindStringSubmatch(strings.ToUpper(sizeStr))
	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid size format, use e.g., 10G or 500M")
	}
	size, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size value: %v", err)
	}
	switch matches[2] {
	case "G":
		return size * 1024 * 1024 * 1024, nil
	case "M":
		return size * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("unsupported unit: %s", matches[2])
	}
}

// allocateMemory 分配指定的字节
func allocateMemory(bytes uint64) ([]byte, error) {
	data := make([]byte, bytes)
	for i := range data {
		data[i] = 1 // Ensure memory is used
	}
	return data, nil
}

// releaseMemory 释放分配的内存
func releaseMemory(data *[]byte) {
	*data = nil
	runtime.GC()
}
