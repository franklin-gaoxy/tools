package cobra_command

import (
	"fmt"
	"many/tools"
	"many/tools/printline"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

// NewParallelCommand creates a command to run tasks in parallel on all target hosts.
func NewParallelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parallel [command]",
		Short: "Run in parallel mode",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hosts, err := GetHosts()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			RunParallel(hosts, strings.Join(args, " "))
		},
	}
	cmd.Flags().StringVarP(&Group, "group", "g", "", "target group (default select all)")
	return cmd
}

// RunParallel executes a command on multiple hosts concurrently.
// It uses a goroutine for each host and waits for all to complete.
func RunParallel(hosts []tools.HostInfo, cmd string) {
	RunParallelWithLimit(hosts, cmd, 0)
}

// RunParallelWithLimit executes a command on multiple hosts concurrently with a limit.
// If limit is 0, it behaves like RunParallel (unlimited concurrency).
func RunParallelWithLimit(hosts []tools.HostInfo, cmd string, limit int) {
	var wg sync.WaitGroup
	var sem chan struct{}

	if limit > 0 {
		sem = make(chan struct{}, limit)
	}

	for _, h := range hosts {
		if limit > 0 {
			sem <- struct{}{}
		}
		wg.Add(1)
		go func(host tools.HostInfo) {
			defer wg.Done()
			if limit > 0 {
				defer func() { <-sem }()
			}
			out, err := tools.RunCommand(host, cmd)
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
