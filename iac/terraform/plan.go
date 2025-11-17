package main

import (
	"context"

	"dagger/terraform/internal/dagger"
)

// Plan génère et affiche un plan d'exécution Terraform
//
// Cette fonction exécute `terraform init` suivi de `terraform plan`.
// Les variables doivent être configurées au préalable via WithVariable().
//
// Exemple:
//
//	dagger call \
//	  with-variable --key vsphere_user --value env:VSPHERE_USER --secret --tf-var \
//	  with-variable --key vsphere_password --value env:VSPHERE_PASSWORD --secret --tf-var \
//	  with-variable --key vsphere_server --value vcenter.local --tf-var \
//	  with-state --backend s3 --bucket my-state --key terraform.tfstate --region us-east-1 \
//	  plan --source ./terraform
func (m *Terraform) Plan(
	ctx context.Context,
	// Répertoire contenant le code Terraform
	source *dagger.Directory,
	// Utiliser -detailed-exitcode (0=no changes, 1=error, 2=changes)
	// +optional
	// +default=false
	detailedExitcode bool,
	// Options supplémentaires pour terraform plan
	// +optional
	planArgs []string,
) (string, error) {
	// Configurer le backend si un état est configuré
	source, err := m.configureBackend(ctx, source)
	if err != nil {
		return "", err
	}

	// Construire le conteneur de base
	container := m.buildContainer(source)

	// Injecter les variables
	container, err = m.injectVariables(ctx, container)
	if err != nil {
		return "", err
	}

	// Exécuter terraform init
	container = container.WithExec([]string{"terraform", "init"})

	// Construire la commande plan
	args := []string{"terraform", "plan"}
	if detailedExitcode {
		args = append(args, "-detailed-exitcode")
	}
	if len(planArgs) > 0 {
		args = append(args, planArgs...)
	}

	// Exécuter terraform plan
	container = container.WithExec(args)

	// Retourner le résultat
	return container.Stdout(ctx)
}
