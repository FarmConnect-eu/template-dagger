package main

// WithState configure le backend Terraform pour la gestion de l'état
// Supporte plusieurs types de backends : s3, gcs, azurerm, local
//
// Le fichier backend.tf est généré dynamiquement au moment de l'exécution
//
// Exemple S3:
//
//	dagger call \
//	  with-state --backend s3 --bucket my-terraform-state --key myapp/terraform.tfstate --region us-east-1 \
//	  plan --source . --workdir terraform
//
// Exemple GCS:
//
//	dagger call \
//	  with-state --backend gcs --bucket my-terraform-state --key myapp/terraform.tfstate \
//	  plan --source . --workdir terraform
//
// Exemple Azure:
//
//	dagger call \
//	  with-state --backend azurerm --bucket mycontainer --key myapp.tfstate \
//	  plan --source . --workdir terraform
//
// Exemple Local:
//
//	dagger call \
//	  with-state --backend local --key terraform.tfstate \
//	  plan --source . --workdir terraform
func (m *Terraform) WithState(
	// Type de backend (s3, gcs, azurerm, local)
	backend string,
	// Nom du bucket/container (non utilisé pour local)
	// +optional
	bucket string,
	// Chemin de la clé/fichier d'état
	// +optional
	key string,
	// Région (utilisé pour S3)
	// +optional
	region string,
) *Terraform {
	stateConfig := &StateConfig{
		Backend: backend,
		Bucket:  bucket,
		Key:     key,
		Region:  region,
	}

	// Deep copy pour éviter les mutations (pattern immutable)
	newVariables := make([]Variable, len(m.Variables))
	copy(newVariables, m.Variables)

	return &Terraform{
		Variables:        newVariables,
		State:            stateConfig,
		TerraformVersion: m.TerraformVersion,
	}
}
