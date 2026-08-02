# Module Dagger Terraform

Module Dagger réutilisable pour exécuter Terraform dans un conteneur de manière portable et reproductible.

## 🎯 Fonctionnalités

- **Plan**: Génère un plan d'exécution Terraform
- **PlanArtifact**: Produit un plan sauvegardé (`tfplan.bin` + `tfplan.json` + `plan.txt` + `changes`)
- **InitArtifact**: Initialise les modules distants pour un scan Trivy exhaustif
- **Apply**: Applique les changements d'infrastructure, ou rejoue un plan sauvegardé
- **Destroy**: Détruit l'infrastructure gérée
- **Validate**: Valide la configuration Terraform
- **Format**: Formate les fichiers Terraform
- **Output**: Récupère les outputs Terraform au format JSON
- **Gestion d'état**: Support backends S3, GCS, Azure, local
- **Variables sécurisées**: Support natif Dagger pour env:, file:, etc.
- **Fichiers tfvars**: Montage automatique via `WithTfVarsFile`

## 🏗️ Architecture

Ce module utilise une **architecture modulaire** avec un pattern immutable :

```
terraform/
├── main.go                    # Structs et types (60 lignes)
├── helpers.go                 # Fonctions utilitaires partagées
├── plan.go                    # Opérations Plan et PlanArtifact
├── init_artifact.go           # Opération InitArtifact
├── deploy.go                  # Opération Apply
├── destroy.go                 # Opération Destroy
├── validate.go                # Opération Validate
├── format.go                  # Opération Format
├── outputs.go                 # Opération Output
├── with_variable.go           # Gestion des variables
├── with_secret.go             # Gestion des variables secrètes
├── with_tfvars.go             # Montage de fichiers .tfvars
├── with_state.go              # Configuration du backend
├── with_terraform_version.go  # Version de Terraform
└── README.md
```

**Avantages** :
- Code modulaire et maintenable
- Chaque fichier < 100 lignes
- Facilement extensible
- Pattern immutable avec deep copy

## 🚀 Quick Start

### Installation

```bash
# Cloner DAGGER-TEMPLATES
cd /path/to/DAGGER-TEMPLATES/iac/terraform

# Générer le SDK
dagger develop

# Lister les fonctions
dagger functions
```

### Utilisation de Base

```bash
# Valider la configuration
dagger call validate --source /path/to/terraform/project

# Plan simple
dagger call plan --source .

# Plan avec variables
dagger call \
  with-variable --key vsphere_user --value env:VSPHERE_USER --secret --tf-var \
  with-variable --key vsphere_password --value env:VSPHERE_PASSWORD --secret --tf-var \
  with-variable --key vsphere_server --value vcenter.local --tf-var \
  plan --source ./terraform

# Apply avec backend S3
dagger call \
  with-variable --key vsphere_user --value env:VSPHERE_USER --secret --tf-var \
  with-variable --key vsphere_password --value env:VSPHERE_PASSWORD --secret --tf-var \
  with-variable --key AWS_ACCESS_KEY_ID --value env:AWS_ACCESS_KEY_ID --secret \
  with-variable --key AWS_SECRET_ACCESS_KEY --value env:AWS_SECRET_ACCESS_KEY --secret \
  with-state --backend s3 --bucket my-terraform-state --key myapp/terraform.tfstate --region us-east-1 \
  apply --source ./terraform --auto-approve
```

## 📦 Fonctions Disponibles

### Configuration (Chainable)

Ces fonctions retournent un nouveau module Terraform et peuvent être chaînées :

#### WithVariable

Ajoute une variable à injecter dans Terraform.

**Paramètres** :
- `key` : Nom de la variable
- `value` : Valeur (supporte `env:`, `file:`, ou valeur littérale)
- `secret` : Marquer comme secret (défaut: `false`)
- `tf-var` : Ajouter le préfixe `TF_VAR_` (défaut: `false`)

**Support des préfixes Dagger** :
- `env:MY_VAR` : Résout la variable d'environnement MY_VAR
- `file:/path/to/file` : Lit le contenu du fichier
- Valeur littérale : Utilisée telle quelle

**Exemples** :
```bash
# Variable Terraform non secrète
dagger call \
  with-variable --key target_node --value pve-node-01 --tf-var \
  plan --source .

# Variable Terraform secrète (depuis env)
dagger call \
  with-variable --key proxmox_api_url --value env:PM_API_URL --secret --tf-var \
  plan --source .

# Variable d'environnement simple (non TF_VAR)
dagger call \
  with-variable --key AWS_ACCESS_KEY_ID --value env:AWS_ACCESS_KEY_ID --secret \
  plan --source .

# Valeur depuis un fichier
dagger call \
  with-variable --key ssh_private_key --value file:/path/to/key --secret --tf-var \
  plan --source .
```

#### WithState

Configure le backend Terraform pour la gestion de l'état.

**Backends supportés** : `s3`, `gcs`, `azurerm`, `local`

**Paramètres** :
- `backend` : Type de backend
- `bucket` : Nom du bucket/container (non utilisé pour local)
- `key` : Chemin de la clé/fichier d'état
- `region` : Région (utilisé pour S3)

**Exemples** :
```bash
# Backend S3
dagger call \
  with-state --backend s3 --bucket my-terraform-state --key myapp/terraform.tfstate --region us-east-1 \
  plan --source .

# Backend GCS
dagger call \
  with-state --backend gcs --bucket my-terraform-state --key myapp/terraform.tfstate \
  plan --source .

# Backend Azure
dagger call \
  with-state --backend azurerm --bucket mycontainer --key myapp.tfstate \
  plan --source .

# Backend Local
dagger call \
  with-state --backend local --key terraform.tfstate \
  plan --source .
```

#### WithTerraformVersion

Configure la version de Terraform à utiliser (défaut: `1.10.6`).

**Exemple** :
```bash
dagger call \
  with-terraform-version --version 1.10.0 \
  plan --source .
```

#### WithTfVarsFile

Monte un fichier `.tfvars`. Les fichiers sont montés sous le nom
`dagger-<n>.auto.tfvars` pour être chargés automatiquement par OpenTofu.
Chaîner l'appel plusieurs fois pour monter plusieurs fichiers.

**Exemple** :
```bash
dagger call \
  with-tf-vars-file --file ./terraform.tfvars \
  plan --source ./terraform
```


### Opérations Terraform

#### Plan

Génère et affiche un plan d'exécution Terraform.

**Paramètres** :
- `source` : Répertoire contenant le code Terraform
- `detailed-exitcode` : Utiliser `-detailed-exitcode` (défaut: `false`)
- `destroy` : Générer un plan de destruction (`tofu plan -destroy`, défaut: `false`)
- `plan-args` : Arguments supplémentaires pour `terraform plan`

**Exemple complet** :
```bash
dagger call \
  with-variable --key vsphere_user --value env:VSPHERE_USER --secret --tf-var \
  with-variable --key vsphere_password --value env:VSPHERE_PASSWORD --secret --tf-var \
  with-variable --key vsphere_server --value vcenter.example.com --tf-var \
  with-variable --key vsphere_datacenter --value DC1 --tf-var \
  with-variable --key AWS_ACCESS_KEY_ID --value env:AWS_ACCESS_KEY_ID --secret \
  with-variable --key AWS_SECRET_ACCESS_KEY --value env:AWS_SECRET_ACCESS_KEY --secret \
  with-state --backend s3 --bucket my-state --key terraform.tfstate --region us-east-1 \
  plan --source ./terraform --detailed-exitcode
```

#### PlanArtifact

Génère un plan **sauvegardé et réutilisable**, et retourne un répertoire contenant :

| Fichier | Contenu |
|---------|---------|
| `tfplan.bin` | Plan binaire, consommable tel quel par `apply --plan-file` |
| `tfplan.json` | `tofu show -json` du plan — entrée du scan `security/trivy scan-plan` |
| `plan.txt` | `tofu show` lisible pour les reviewers (valeurs sensibles masquées) |
| `changes` | `"true"` si des changements sont en attente, sinon `"false"` |

L'appel n'échoue pas quand des changements existent (code de sortie 2 intercepté) :
la barrière CI se décide en lisant `changes`.

**Paramètres** :
- `source` : Répertoire contenant le code Terraform
- `subpath` : Sous-chemin relatif dans `source` (défaut: `.`)
- `destroy` : Générer un plan de destruction (défaut: `false`)
- `plan-args` : Arguments supplémentaires pour `tofu plan`

**Exemple — plan → scan → apply du plan exact** :
```bash
# 1. Produire l'artefact de plan
dagger call \
  with-secret --key proxmox_api_token_secret --value env:PM_TOKEN_SECRET --tf-var \
  with-state --backend s3 --bucket terraform --key infra-nfs/state.tfstate --region main \
  plan-artifact --source . --subpath terraform \
  export --path /tmp/plan

# 2. Scanner les misconfigurations du plan
dagger call -m ../../security/trivy \
  scan-plan --plan /tmp/plan/tfplan.json --fail-on-findings

# 3. Appliquer exactement ce plan (pas de recalcul)
dagger call \
  with-secret --key proxmox_api_token_secret --value env:PM_TOKEN_SECRET --tf-var \
  with-state --backend s3 --bucket terraform --key infra-nfs/state.tfstate --region main \
  apply --source . --subpath terraform --plan-file /tmp/plan/tfplan.bin
```

#### InitArtifact

Exécute `tofu init -backend=false` et retourne le répertoire de travail initialisé,
`.terraform/modules` inclus. Destiné à être scanné par `security/trivy scan-config` :
Trivy voit alors le contenu réel des modules distants, et pas seulement le code wrapper.
`.terraform/providers` est retiré de l'artefact (inutile à `trivy config`, et volumineux).

**Exemple** :
```bash
dagger call init-artifact --source . --subpath terraform export --path /tmp/tf-init
dagger call -m ../../security/trivy scan-config --source /tmp/tf-init --fail-on-findings
```

#### Apply

Applique les changements Terraform à l'infrastructure.

Sans `--plan-file`, le plan est recalculé au moment de l'apply. Avec `--plan-file`,
le plan produit par `PlanArtifact` est appliqué **tel quel** : si l'état a changé depuis
sa génération, OpenTofu refuse le plan obsolète (« saved plan is stale »).

> Les credentials providers restent nécessaires même avec un plan sauvegardé : ils ne
> sont pas stockés dans le plan. Continuer à chaîner `with-secret` / `with-variable`.

**Paramètres** :
- `source` : Répertoire contenant le code Terraform
- `plan-file` : Plan sauvegardé (`tfplan.bin`) à appliquer tel quel (optionnel)
- `apply-args` : Arguments supplémentaires pour `terraform apply`

**Exemple** :
```bash
dagger call \
  with-variable --key vsphere_user --value env:VSPHERE_USER --secret --tf-var \
  with-variable --key vsphere_password --value env:VSPHERE_PASSWORD --secret --tf-var \
  with-variable --key AWS_ACCESS_KEY_ID --value env:AWS_ACCESS_KEY_ID --secret \
  with-variable --key AWS_SECRET_ACCESS_KEY --value env:AWS_SECRET_ACCESS_KEY --secret \
  with-state --backend s3 --bucket my-state --key terraform.tfstate --region us-east-1 \
  apply --source ./terraform --auto-approve
```

#### Destroy

Détruit l'infrastructure gérée par Terraform.

**Paramètres** :
- `source` : Répertoire contenant le code Terraform
- `auto-approve` : Détruire sans confirmation (défaut: `false`)
- `destroy-args` : Arguments supplémentaires pour `terraform destroy`

**Exemple** :
```bash
dagger call \
  with-variable --key vsphere_user --value env:VSPHERE_USER --secret --tf-var \
  with-variable --key vsphere_password --value env:VSPHERE_PASSWORD --secret --tf-var \
  with-variable --key AWS_ACCESS_KEY_ID --value env:AWS_ACCESS_KEY_ID --secret \
  with-variable --key AWS_SECRET_ACCESS_KEY --value env:AWS_SECRET_ACCESS_KEY --secret \
  with-state --backend s3 --bucket my-state --key terraform.tfstate --region us-east-1 \
  destroy --source ./terraform --auto-approve
```

#### Validate

Valide la syntaxe et la configuration Terraform.

**Paramètres** :
- `source` : Répertoire contenant le code Terraform

**Exemple** :
```bash
dagger call validate --source ./terraform
```

#### Format

Formate les fichiers Terraform selon les standards.

**Paramètres** :
- `source` : Répertoire contenant le code Terraform
- `check` : Vérifier seulement sans modifier (défaut: `false`)
- `recursive` : Formatter récursivement (défaut: `true`)

**Exemples** :
```bash
# Formatter les fichiers
dagger call format --source ./terraform

# Vérifier le formatage (CI)
dagger call format --source ./terraform --check
```

#### Output

Récupère les outputs Terraform.

**Paramètres** :
- `source` : Répertoire contenant le code Terraform
- `output-name` : Nom d'un output spécifique (vide = tous)
- `as-json` : Format JSON (défaut: `true`)

**Exemples** :
```bash
# Tous les outputs
dagger call \
  with-variable --key AWS_ACCESS_KEY_ID --value env:AWS_ACCESS_KEY_ID --secret \
  with-variable --key AWS_SECRET_ACCESS_KEY --value env:AWS_SECRET_ACCESS_KEY --secret \
  with-state --backend s3 --bucket my-state --key terraform.tfstate --region us-east-1 \
  output --source ./terraform

# Output spécifique
dagger call \
  with-variable --key AWS_ACCESS_KEY_ID --value env:AWS_ACCESS_KEY_ID --secret \
  with-variable --key AWS_SECRET_ACCESS_KEY --value env:AWS_SECRET_ACCESS_KEY --secret \
  with-state --backend s3 --bucket my-state --key terraform.tfstate --region us-east-1 \
  output --source ./terraform --output-name vm_ip_addresses
```

## 🔧 Exemples d'Utilisation

### Scénario 1 : vSphere + S3 Backend

```bash
dagger call \
  with-terraform-version --version 1.9.8 \
  with-variable --key vsphere_user --value env:VSPHERE_USER --secret --tf-var \
  with-variable --key vsphere_password --value env:VSPHERE_PASSWORD --secret --tf-var \
  with-variable --key vsphere_server --value vcenter.example.com --tf-var \
  with-variable --key vsphere_datacenter --value DC1 --tf-var \
  with-variable --key vsphere_network --value PROD-DMZ --tf-var \
  with-variable --key AWS_ACCESS_KEY_ID --value env:AWS_ACCESS_KEY_ID --secret \
  with-variable --key AWS_SECRET_ACCESS_KEY --value env:AWS_SECRET_ACCESS_KEY --secret \
  with-state --backend s3 --bucket my-terraform-state --key vsphere/terraform.tfstate --region us-east-1 \
  plan --source /path/to/vsphere/project
```

### Scénario 2 : AWS + Backend Local

```bash
dagger call \
  with-variable --key AWS_ACCESS_KEY_ID --value env:AWS_ACCESS_KEY_ID --secret \
  with-variable --key AWS_SECRET_ACCESS_KEY --value env:AWS_SECRET_ACCESS_KEY --secret \
  with-variable --key aws_region --value us-east-1 --tf-var \
  with-variable --key instance_type --value t3.micro --tf-var \
  with-state --backend local --key terraform.tfstate \
  apply --source ./terraform --auto-approve
```

### Scénario 3 : FortiGate + GCS Backend

```bash
dagger call \
  with-variable --key fortigate_hostname --value env:FORTIGATE_HOSTNAME --secret --tf-var \
  with-variable --key fortigate_token --value env:FORTIGATE_TOKEN --secret --tf-var \
  with-variable --key fortigate_port --value 443 --tf-var \
  with-variable --key GOOGLE_CREDENTIALS --value env:GOOGLE_CREDENTIALS --secret \
  with-state --backend gcs --bucket my-terraform-state --key fortigate/terraform.tfstate \
  plan --source ./terraform
```

### Scénario 4 : Pipeline CI/CD

```bash
# Étape 1 : Validate
dagger call validate --source ./terraform

# Étape 2 : Format check
dagger call format --source ./terraform --check

# Étape 3 : Plan
dagger call \
  with-variable --key api_key --value env:API_KEY --secret --tf-var \
  with-variable --key AWS_ACCESS_KEY_ID --value env:AWS_ACCESS_KEY_ID --secret \
  with-variable --key AWS_SECRET_ACCESS_KEY --value env:AWS_SECRET_ACCESS_KEY --secret \
  with-state --backend s3 --bucket ci-terraform-state --key app/terraform.tfstate --region us-east-1 \
  plan --source ./terraform --detailed-exitcode

# Étape 4 : Apply (si plan ok)
dagger call \
  with-variable --key api_key --value env:API_KEY --secret --tf-var \
  with-variable --key AWS_ACCESS_KEY_ID --value env:AWS_ACCESS_KEY_ID --secret \
  with-variable --key AWS_SECRET_ACCESS_KEY --value env:AWS_SECRET_ACCESS_KEY --secret \
  with-state --backend s3 --bucket ci-terraform-state --key app/terraform.tfstate --region us-east-1 \
  apply --source ./terraform --auto-approve
```

## 🔐 Sécurité

### Gestion des Secrets

Les secrets ne sont **jamais** exposés dans les logs grâce à l'utilisation de `--secret` :

```bash
# ✅ CORRECT : Secret protégé
dagger call \
  with-variable --key password --value env:MY_PASSWORD --secret --tf-var \
  plan --source .

# ❌ INCORRECT : Valeur visible dans les logs
dagger call \
  with-variable --key password --value my-secret-password --tf-var \
  plan --source .
```

### Variables d'Environnement Backend

Pour les backends S3/GCS/Azure, les credentials sont injectés comme variables d'environnement :

```bash
# Backend S3 (avec credentials AWS)
dagger call \
  with-variable --key AWS_ACCESS_KEY_ID --value env:AWS_ACCESS_KEY_ID --secret \
  with-variable --key AWS_SECRET_ACCESS_KEY --value env:AWS_SECRET_ACCESS_KEY --secret \
  with-state --backend s3 --bucket my-state --key terraform.tfstate --region us-east-1 \
  plan --source ./terraform

# Backend GCS (avec credentials Google Cloud)
dagger call \
  with-variable --key GOOGLE_CREDENTIALS --value env:GOOGLE_CREDENTIALS --secret \
  with-state --backend gcs --bucket my-state --key terraform.tfstate \
  plan --source ./terraform

# Backend Azure (avec credentials Azure)
dagger call \
  with-variable --key ARM_CLIENT_ID --value env:ARM_CLIENT_ID --secret \
  with-variable --key ARM_CLIENT_SECRET --value env:ARM_CLIENT_SECRET --secret \
  with-variable --key ARM_TENANT_ID --value env:ARM_TENANT_ID --secret \
  with-variable --key ARM_SUBSCRIPTION_ID --value env:ARM_SUBSCRIPTION_ID --secret \
  with-state --backend azurerm --bucket mycontainer --key terraform.tfstate \
  plan --source ./terraform
```

## 🧪 Tests

```bash
# Tester que le module est chargé
dagger call test

# Valider la configuration
dagger call validate --source ./examples/simple

# Vérifier le formatage
dagger call format --source ./examples/simple --check
```

## 📚 Ressources

- [Documentation Terraform](https://www.terraform.io/docs)
- [Documentation Dagger](https://docs.dagger.io)
- [Terraform Backends](https://www.terraform.io/docs/language/settings/backends)
- [Dagger Secrets](https://docs.dagger.io/api/secrets)

## 🤝 Contribution

Pour ajouter une nouvelle fonctionnalité :

1. Créer un nouveau fichier (ex: `init.go`)
2. Implémenter la fonction sur le type `*Terraform`
3. Utiliser les helpers partagés (`buildContainer`, `injectVariables`, `configureBackend`)
4. Ajouter la documentation dans ce README

Exemple :

```go
// init.go
package main

import (
    "context"
    "dagger/terraform/internal/dagger"
)

func (m *Terraform) Init(ctx context.Context, source *dagger.Directory) (string, error) {
    container := m.buildContainer(source)
    container, _ = m.injectVariables(ctx, container)
    return container.WithExec([]string{"terraform", "init"}).Stdout(ctx)
}
```

## 📄 Licence

MIT
