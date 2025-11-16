// Ansible CI/CD module for Dagger.
// Automates Ansible workflows with secret management.

package main

// Variable to inject into Ansible (supports literal, env://, file://)
type Variable struct {
	Key      string
	Value    string
	IsSecret bool
}

// Ansible module
type Ansible struct {
	Variables      []Variable
	AnsibleVersion string
}
