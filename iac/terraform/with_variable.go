package main

// WithVariable ajoute une variable au module Terraform
// Supporte les valeurs littérales et les préfixes Dagger (env://, file://, etc.)
//
// Les préfixes sont gérés nativement par Dagger :
//   - env:MY_VAR : Résout la variable d'environnement MY_VAR
//   - file:/path/to/file : Lit le contenu du fichier
//   - Valeur littérale : Utilisée telle quelle
//
// Exemple:
//
//	dagger call \
//	  with-variable --key proxmox_api_url --value env:PM_API_URL --secret --tf-var \
//	  with-variable --key target_node --value pve-node-01 --tf-var \
//	  plan --source . --workdir terraform
func (m *Terraform) WithVariable(
	// Nom de la variable
	key string,
	// Valeur (supporte env:, file:, ou valeur littérale)
	value string,
	// Marquer comme secret (injecté via WithSecretVariable)
	// +optional
	// +default=false
	secret bool,
	// Ajouter le préfixe TF_VAR_ (pour variables Terraform)
	// +optional
	// +default=false
	tfVar bool,
) *Terraform {
	newVar := Variable{
		Key:      key,
		Value:    value,
		IsSecret: secret,
		TfVar:    tfVar,
	}

	// Deep copy pour éviter les mutations (pattern immutable)
	newVariables := make([]Variable, len(m.Variables), len(m.Variables)+1)
	copy(newVariables, m.Variables)

	return &Terraform{
		Variables:        append(newVariables, newVar),
		State:            m.State,
		TerraformVersion: m.TerraformVersion,
	}
}
