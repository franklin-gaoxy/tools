package tools

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func BoundBobraAgrs() *cobra.Command {
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

	// 指定使用的文件
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "./nodelist", "hosts file")
	// 并行
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
	// 串行
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
	// 批次
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
	// 每批次的主机
	batchCmd.Flags().IntVarP(&batchSize, "number", "n", 5, "batch size")
	// 绑定
	rootCmd.AddCommand(parallelCmd, serialCmd, batchCmd)

	// hello 测试连接
	helloCmd := &cobra.Command{
		Use:   "test connect remote host",
		Short: "Run test connect",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			hosts, err := readHosts(file)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			TestConnection(hosts)
		},
	}
	rootCmd.AddCommand(helloCmd)

	// version
	versionCmd := &cobra.Command{
		Use:   "show version",
		Short: "show version",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("v0.0.1")
		},
	}
	rootCmd.AddCommand(versionCmd)

	return rootCmd
}

func ExecuteStart(rootCmd *cobra.Command) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
