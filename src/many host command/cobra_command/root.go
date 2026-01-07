package cobra_command

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var (
	File      string // Path to the host list file
	Verbosity int    // Log verbosity level
	Group     string // Target host group to filter
)

// NewRootCommand creates the root Cobra command for the ssh-tool application.
// It sets up global flags and the default run behavior (Parallel execution).
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ssh-tool [command]",
		Short: "SSH execution tool",
		Args:  cobra.MinimumNArgs(1), // 至少需要一个命令
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			klog.InitFlags(nil)
			flag.Set("v", fmt.Sprintf("%d", Verbosity))
			flag.Set("logtostderr", "true")
		},
		Run: func(cmd *cobra.Command, args []string) {
			hosts, err := GetHosts()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			RunParallel(hosts, strings.Join(args, " "))
		},
	}

	// 指定使用的文件
	rootCmd.PersistentFlags().StringVarP(&File, "file", "f", "./nodelist", "hosts file")
	// 指定日志级别
	rootCmd.PersistentFlags().IntVarP(&Verbosity, "verbose", "v", 0, "log level")

	// 为了让 root command 也能使用 group 过滤（因为它默认行为是 parallel），我们也可以在这里添加
	// 但是用户要求给 batch, serial, parallel 增加，root 通常行为跟随 parallel
	// 如果用户直接运行 ssh-tool -g group cmd，也应该生效
	rootCmd.Flags().StringVarP(&Group, "group", "g", "", "target group (default select all)")

	rootCmd.AddCommand(NewParallelCommand())
	rootCmd.AddCommand(NewSerialCommand())
	rootCmd.AddCommand(NewBatchCommand())
	rootCmd.AddCommand(NewScpCommand())
	rootCmd.AddCommand(NewTestCommand())
	rootCmd.AddCommand(NewVersionCommand())

	return rootCmd
}

// Execute is the main entry point for the application.
// It executes the root command and handles any errors.
func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
