package cobra_command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCommand creates a command to display the application version.
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "show version",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("v0.0.1")
		},
	}
}
