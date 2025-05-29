package src

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"math"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// RunCPULoad: run cpu load
func RunCPULoad(cmd *cobra.Command, args []string) {
	load, err := cmd.Flags().GetString("load")
	RunTime, err := cmd.Flags().GetInt("time")
	if err != nil {
		fmt.Printf("Failed to obtain parameters!\n%s", err)
		os.Exit(1)
	}

	// 类型转换
	targetUsage, err := strconv.ParseFloat(load, 64)
	if err != nil {
		fmt.Printf("Please ensure that the correct parameters are passed in!\nSupports both int and float64,Please enter an integer such as 1 or a decimal such as 0.8\n")
		fmt.Println(err)
		os.Exit(1)
	}

	if targetUsage >= 1 {
		fmt.Printf("full load operation, run time: %sm\n", RunTime)
		RunFullCPU(RunTime)
	} else if targetUsage < 1 {
		fmt.Printf("percentage load operation, run time: %sm\n", RunTime)
		RunPercentageCPU(targetUsage, RunTime)
	}

}

/*
tools: Public code area
*/
// 小块计算任务（避免长时间无休眠）
func busyWorkChunk() float64 {
	var sum float64
	for i := 0; i < 1_000; i++ {
		sum += math.Pow(math.Pi, float64(i%50))
	}
	return sum // 返回结果防止被优化
}

// RunPercentageCPU 按百分比方式运行 targetUsage:CPU目标使用率 如 0.8
func RunPercentageCPU(targetUsage float64, RunTime int) {
	// 检查使用率参数是否正确
	if targetUsage < 0.01 || targetUsage > 0.99 {
		fmt.Printf("Invalid usage rate: %.2f (Must be between 0.01-0.99)\n", targetUsage)
		return
	}

	// 确保在函数退出时取消所有goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	// 计算需要休眠的时间 以确保能控制CPU使用率
	workDuration := time.Duration(float64(time.Second) * targetUsage / float64(numCPU))
	sleepDuration := time.Duration(float64(time.Second) * (1 - targetUsage) / float64(numCPU))

	// 为每个CPU核心启动一个goroutine
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					// 上下文取消时退出
					return
				default:
					// 获取开始时间
					start := time.Now()
					for time.Since(start) < workDuration {
						busyWorkChunk()
					}

					// 睡眠 停止计算
					time.Sleep(sleepDuration)
				}
			}
		}(i)
	}

	// 如果指定了退出时间 根据运行时间参数处理退出
	if RunTime > 0 {

		// 创建定时器
		timeout := time.Duration(RunTime) * time.Minute
		timer := time.NewTimer(timeout)
		defer timer.Stop()

		// 创建worker完成通道
		workerDone := make(chan struct{})
		go func() {
			wg.Wait()
			close(workerDone)
		}()

		// 等待超时或所有worker提前完成
		select {
		case <-workerDone:
			fmt.Println("worker complete.")
		case <-timer.C:
			// 到达退出时间 发送取消信号
			cancel()
			// 等待goroutine全部退出
			wg.Wait()
		}
	} else {
		// 一直运行
		wg.Wait()
	}
}

// RunFullCPU 将CPU跑到100%
func RunFullCPU(RunTime int) {
	startTime := time.Now()
	// 确保函数退出时清理资源
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 检查是否到达了退出时间 每分钟检查一次
			ticker := time.NewTicker(time.Minute)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return // 上下文取消时立即退出
				case <-ticker.C:
					if RunTime > 0 && time.Since(startTime) >= time.Duration(RunTime)*time.Minute {
						return // 达到指定时间退出
					}
				default:
					busyWorkChunk()
				}
			}
		}(i)
	}

	if RunTime > 0 {
		// 创建worker完成通知通道
		workerDone := make(chan struct{})
		go func() {
			wg.Wait()
			close(workerDone)
		}()

		// 等待超时或所有worker提前完成
		select {
		case <-workerDone:
			fmt.Println("All workers completed")
		case <-time.After(time.Duration(RunTime) * time.Minute):
			fmt.Println("Timeout reached, canceling workers")
			cancel()  // 通知所有worker退出
			wg.Wait() // 等待所有worker实际退出
		}
	} else {
		fmt.Println("Running indefinitely. Press Ctrl+C to terminate.")
		wg.Wait() // 永久阻塞
	}
}
