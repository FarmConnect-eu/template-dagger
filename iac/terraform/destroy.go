package main

import (
	"context"

	"dagger/terraform/internal/dagger"
)

// Destroy détruit l'infrastructure gérée par Terraform
//
// Cette fonction exécute `terraform init` suivi de `terraform destroy`.
// Les variables doivent être configurées au préalable via WithVariable().
//
// Exemple:
//
//	dagger call \
//	  with-variable --key vsphere_user --value env:VSPHERE_USER --secret --tf-var \
//	  with-variable --key vsphere_password --value env:VSPHERE_PASSWORD --secret --tf-var \
//	  with-state --backend s3 --bucket my-state --key terraform.tfstate --region us-east-1 \
//	  destroy --source ./terraform --auto-approve
func (m *Terraform) Destroy(
	ctx context.Context,
	// Répertoire contenant le code Terraform
	source *dagger.Directory,
	// Détruire automatiquement sans confirmation
	// +optional
	// +default=false
	autoApprove bool,
	// Options supplémentaires pour terraform destroy
	// +optional
	destroyArgs []string,
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

	// Construire la commande destroy
	args := []string{"terraform", "destroy"}
	if autoApprove {
		args = append(args, "-auto-approve")
	}
	if len(destroyArgs) > 0 {
		args = append(args, destroyArgs...)
	}

	// Exécuter terraform destroy
	container = container.WithExec(args)

	// Retourner le résultat
	return container.Stdout(ctx)
}
