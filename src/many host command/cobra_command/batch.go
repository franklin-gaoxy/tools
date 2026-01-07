package cobra_command

import (
	"fmt"
	"many/tools"
	"strings"

	"github.com/spf13/cobra"
)

// NewBatchCommand creates a command to run tasks in batches.
// It supports a --number flag to specify the batch size.
func NewBatchCommand() *cobra.Command {
	var batchSize int
	cmd := &cobra.Command{
		Use:   "batch [command]",
		Short: "Run in batch mode",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hosts, err := tools.ReadHosts(File, Group)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			runBatch(hosts, strings.Join(args, " "), batchSize)
		},
	}
	cmd.Flags().IntVarP(&batchSize, "number", "n", 5, "batch size")
	cmd.Flags().StringVarP(&Group, "group", "g", "", "target group")
	return cmd
}

// runBatch executes a command on hosts in batches.
// It processes 'batchSize' hosts concurrently, then moves to the next batch.
func runBatch(hosts []tools.HostInfo, cmd string, batchSize int) {
	for i := 0; i < len(hosts); i += batchSize {
		end := i + batchSize
		if end > len(hosts) {
			end = len(hosts)
		}
		batch := hosts[i:end]
		tools.RunParallel(batch, cmd)
	}
}
