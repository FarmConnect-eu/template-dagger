package main

import (
	"context"

	"dagger/terraform/internal/dagger"
)

// Apply applique les changements Terraform à l'infrastructure
//
// Cette fonction exécute `terraform init` suivi de `terraform apply`.
// Les variables doivent être configurées au préalable via WithVariable().
//
// Exemple:
//
//	dagger call \
//	  with-variable --key vsphere_user --value env:VSPHERE_USER --secret --tf-var \
//	  with-variable --key vsphere_password --value env:VSPHERE_PASSWORD --secret --tf-var \
//	  with-state --backend s3 --bucket my-state --key terraform.tfstate --region us-east-1 \
//	  apply --source ./terraform --auto-approve
func (m *Terraform) Apply(
	ctx context.Context,
	// Répertoire contenant le code Terraform
	source *dagger.Directory,
	// Appliquer automatiquement sans confirmation
	// +optional
	// +default=false
	autoApprove bool,
	// Options supplémentaires pour terraform apply
	// +optional
	applyArgs []string,
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

	// Construire la commande apply
	args := []string{"terraform", "apply"}
	if autoApprove {
		args = append(args, "-auto-approve")
	}
	if len(applyArgs) > 0 {
		args = append(args, applyArgs...)
	}

	// Exécuter terraform apply
	container = container.WithExec(args)

	// Retourner le résultat
	return container.Stdout(ctx)
}
