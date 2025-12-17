// Package main provides a Dagger module for executing Ansible in containers.
package main

import (
	"context"
	_ "embed"

	"dagger/ansible/internal/dagger"
)

var (
	//go:embed config/requirements.txt
	requirementsTxt string

	//go:embed config/default_vars.yml
	defaultVarsYml string

	//go:embed config/ansible.cfg
	ansibleCfg string
)

type Variable struct {
	Key         string
	Value       string
	SecretValue *dagger.Secret
}

type KeyValue struct {
	Key   string
	Value string
}

type Ansible struct {
	Variables      []Variable
	AnsibleVersion string
	Inventory      *dagger.File
	Requirements   *dagger.File
	RolesPath      *dagger.Directory
	Templates      *dagger.Directory
	GroupVars      *dagger.Directory
	ExtraVars      []KeyValue
	Tags           []string
	SkipTags       []string
}

func New() *Ansible {
	return &Ansible{
		Variables:      []Variable{},
		AnsibleVersion: "11.1.0",
		Inventory:      nil,
		Requirements:   nil,
		RolesPath:      nil,
		Templates:      nil,
		GroupVars:      nil,
		ExtraVars:      []KeyValue{},
		Tags:           []string{},
		SkipTags:       []string{},
	}
}

// Test verifies the module loads correctly.
func (m *Ansible) Test(ctx context.Context) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"echo", "Ansible module loaded successfully"}).
		Stdout(ctx)
}
