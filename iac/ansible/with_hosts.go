package main

import (
	"fmt"
	"strings"
)

// WithHosts creates a dynamic inventory from a comma-separated list of hosts.
// This is useful when you have a list of IPs from Terraform outputs.
// Example: with-hosts --hosts "192.168.1.10,192.168.1.11" --user "admincd24"
func (m *Ansible) WithHosts(
	// Comma-separated list of hosts (e.g., "192.168.1.10,192.168.1.11")
	hosts string,
	// Group name for the hosts (default: "all")
	// +optional
	// +default="all"
	group string,
	// SSH user for connecting to hosts (sets ansible_user)
	// +optional
	user string,
) *Ansible {
	if group == "" {
		group = "all"
	}

	hostList := strings.Split(hosts, ",")
	var inventoryContent strings.Builder

	inventoryContent.WriteString(fmt.Sprintf("[%s]\n", group))
	for _, host := range hostList {
		host = strings.TrimSpace(host)
		if host != "" {
			inventoryContent.WriteString(host + "\n")
		}
	}

	// Add group variables if user is specified
	if user != "" {
		inventoryContent.WriteString(fmt.Sprintf("\n[%s:vars]\n", group))
		inventoryContent.WriteString(fmt.Sprintf("ansible_user=%s\n", user))
	}

	inventoryFile := dag.Directory().
		WithNewFile("inventory", inventoryContent.String()).
		File("inventory")

	newVariables := make([]Variable, len(m.Variables))
	copy(newVariables, m.Variables)

	newExtraVars := make([]KeyValue, len(m.ExtraVars))
	copy(newExtraVars, m.ExtraVars)

	newTags := make([]string, len(m.Tags))
	copy(newTags, m.Tags)

	newSkipTags := make([]string, len(m.SkipTags))
	copy(newSkipTags, m.SkipTags)

	return &Ansible{
		Variables:      newVariables,
		AnsibleVersion: m.AnsibleVersion,
		Inventory:      inventoryFile,
		Requirements:   m.Requirements,
		RolesPath:      m.RolesPath,
		Templates:      m.Templates,
		GroupVars:      m.GroupVars,
		ExtraVars:      newExtraVars,
		Tags:           newTags,
		SkipTags:       newSkipTags,
	}
}
