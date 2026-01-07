package cobra_command

import (
	"fmt"
	"many/tools"
)

// GetHosts reads the host list file and filters by group if specified.
// It uses the global File and Group variables.
func GetHosts() ([]tools.HostInfo, error) {
	hosts, err := tools.ReadHosts(File, Group)
	if err != nil {
		return nil, fmt.Errorf("failed to read hosts: %w", err)
	}
	return hosts, nil
}
