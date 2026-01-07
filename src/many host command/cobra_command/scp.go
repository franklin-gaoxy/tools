package cobra_command

import (
	"fmt"
	"many/tools"

	"github.com/spf13/cobra"
)

// NewScpCommand creates a command to copy files to target hosts.
func NewScpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scp [src] [dest]",
		Short: "Copy files to remote hosts",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			hosts, err := GetHosts()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			src := args[0]
			dest := args[1]
			tools.RunScp(hosts, src, dest)
		},
	}
	cmd.Flags().StringVarP(&Group, "group", "g", "", "target group")
	return cmd
}
