// Module Dagger pour exécuter Terraform dans un conteneur
// Ce module fournit des fonctions réutilisables pour gérer l'infrastructure avec Terraform
// de manière portable et reproductible.
//
// Le module utilise un pattern immutable avec accumulation d'état :
//
//	dagger call \
//	  with-variable --key vsphere_user --value env:VSPHERE_USER --secret --tf-var \
//	  with-variable --key vsphere_password --value env:VSPHERE_PASSWORD --secret --tf-var \
//	  with-variable --key vsphere_server --value vcenter.local --tf-var \
//	  with-state --backend s3 --bucket my-state --key terraform.tfstate --region us-east-1 \
//	  plan --source . --workdir terraform
//
package main

import (
	"context"
)

// Variable représente une variable à injecter dans Terraform
// Supporte les valeurs littérales et les préfixes Dagger (env://, file://, etc.)
type Variable struct {
	Key      string // Nom de la variable
	Value    string // Valeur (supporte env://, file://, etc.)
	IsSecret bool   // Si true, injecté comme secret
	TfVar    bool   // Si true, préfixe avec TF_VAR_
}

// StateConfig configure le backend Terraform
type StateConfig struct {
	Backend string // Type de backend (s3, gcs, azurerm, local)
	Bucket  string // Nom du bucket/container
	Key     string // Chemin de la clé/fichier d'état
	Region  string // Région (pour S3)
}

// Terraform est le module principal
type Terraform struct {
	Variables        []Variable    // Variables accumulées
	State            *StateConfig  // Configuration du backend
	TerraformVersion string        // Version de Terraform à utiliser
}

// New crée une nouvelle instance du module Terraform
func New() *Terraform {
	return &Terraform{
		Variables:        []Variable{},
		State:            nil,
		TerraformVersion: "1.9.8",
	}
}

// Test vérifie que le module est chargé correctement
func (m *Terraform) Test(ctx context.Context) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"echo", "Terraform module loaded successfully"}).
		Stdout(ctx)
}
