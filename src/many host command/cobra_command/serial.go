package cobra_command

import (
	"fmt"
	"many/tools"
	"strings"

	"github.com/spf13/cobra"
)

// NewSerialCommand creates a command to run tasks sequentially on target hosts.
func NewSerialCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serial [command]",
		Short: "Run in serial mode",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hosts, err := GetHosts()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			runSerial(hosts, strings.Join(args, " "))
		},
	}
	cmd.Flags().StringVarP(&Group, "group", "g", "", "target group (default select all)")
	return cmd
}

// runSerial executes a command on multiple hosts sequentially.
// It processes one host at a time, which is useful for debugging or avoiding high load.
func runSerial(hosts []tools.HostInfo, cmd string) {
	for _, h := range hosts {
		out, err := tools.RunCommand(h, cmd)
		header := fmt.Sprintf("\n[%s]\n", h.IP)
		if err != nil {
			fmt.Printf("%s%s", header, out)
			fmt.Printf("[%s] Command failed: %v\n", h.IP, err)
		} else {
			fmt.Printf("%s%s", header, out)
		}
	}
}
