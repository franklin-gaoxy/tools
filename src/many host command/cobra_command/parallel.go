package cobra_command

import (
	"fmt"
	"many/tools"
	"strings"

	"github.com/spf13/cobra"
)

// NewParallelCommand creates a command to run tasks in parallel on all target hosts.
func NewParallelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parallel [command]",
		Short: "Run in parallel mode",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hosts, err := tools.ReadHosts(File, Group)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			tools.RunParallel(hosts, strings.Join(args, " "))
		},
	}
	cmd.Flags().StringVarP(&Group, "group", "g", "", "target group")
	return cmd
}
