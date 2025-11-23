# Module Dagger Ansible

Module Dagger réutilisable pour exécuter Ansible dans un conteneur de manière portable et reproductible.

## 🎯 Fonctionnalités

- **RunPlaybook**: Exécute des playbooks Ansible
- **Gestion d'inventaire**: Support hosts dynamiques ou fichiers d'inventaire
- **Variables sécurisées**: Support natif Dagger pour env:, file:, etc.
- **Tags**: Support des tags et skip-tags Ansible
- **Extra Vars**: Injection de variables extra
- **Mode Check**: Dry-run sans modifications
- **Secrets Dagger**: Intégration native avec les secrets Dagger (file:, env:)

## 🏗️ Architecture

Ce module utilise une **architecture modulaire** avec un pattern immutable :

```
ansible/
├── main.go                   # Structs et types
├── helpers.go                # Fonctions utilitaires partagées
├── playbook.go               # Exécution de playbooks
├── with_variable.go          # Gestion des variables
├── with_secret.go            # Gestion des secrets
├── with_ansible_version.go   # Version d'Ansible
├── with_inventory.go         # Configuration inventaire
├── config/                   # Fichiers de configuration embarqués
│   ├── requirements.txt
│   └── default_vars.yml
├── scripts/                  # Scripts embarqués
│   └── run-ansible.sh
└── README.md
```

**Avantages** :
- Code modulaire et maintenable
- Pattern immutable avec deep copy
- Support natif des secrets Dagger
- Simplicité et portabilité

## 🚀 Quick Start

### Installation

```bash
# Cloner DAGGER-TEMPLATES
cd /path/to/DAGGER-TEMPLATES/iac/ansible

# Générer le SDK
dagger develop

# Lister les fonctions
dagger functions
```

### Utilisation de Base

```bash
# Exécuter un playbook simple
dagger call \
  with-inventory --hosts "server1.example.com,server2.example.com" \
  with-variable --key ansible_user --value ubuntu \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa

# Avec variables secrètes
dagger call \
  with-inventory --hosts "server1.example.com,server2.example.com" \
  with-variable --key ansible_user --value ubuntu \
  with-secret --key ansible_password --value env:ANSIBLE_PASSWORD \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa
```

## 📦 Fonctions Disponibles

### Configuration (Chainable)

Ces fonctions retournent un nouveau module Ansible et peuvent être chaînées :

#### WithVariable

Ajoute une variable (non-secrète) à injecter dans Ansible.

**Paramètres** :
- `key` : Nom de la variable
- `value` : Valeur de la variable

**Exemple** :
```bash
# Variable simple
dagger call \
  with-variable --key ansible_user --value ubuntu \
  with-variable --key ansible_port --value 22 \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa
```

#### WithSecret

Ajoute un secret au module Ansible. Les secrets sont injectés via WithSecretVariable et ne sont pas exposés dans les logs.

**Paramètres** :
- `key` : Nom de la variable secrète
- `value` : Valeur du secret (supporte `env:`, `file:`, etc.)

**Exemples** :
```bash
# Secret depuis environment
dagger call \
  with-secret --key ansible_password --value env:ANSIBLE_PASSWORD \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa

# Secret depuis un fichier
dagger call \
  with-secret --key api_token --value file:/path/to/token \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa
```

#### WithExtraVar

Ajoute une variable extra pour Ansible (--extra-vars).

**Exemple** :
```bash
dagger call \
  with-extra-var --key deployment_env --value production \
  with-extra-var --key app_version --value 1.2.3 \
  run-playbook --playbook deploy.yml --project . --ssh-private-key file:~/.ssh/id_rsa
```

#### WithInventory

Configure l'inventaire Ansible (hosts ou fichier).

**Paramètres** :
- `hosts` : Liste de hosts séparés par des virgules
- `path` : Chemin vers un fichier d'inventaire

**Exemples** :
```bash
# Hosts dynamiques
dagger call \
  with-inventory --hosts "web1.example.com,web2.example.com,db1.example.com" \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa

# Fichier d'inventaire
dagger call \
  with-inventory --path inventory/production.ini \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa
```

#### WithInventoryVar

Ajoute des variables d'inventaire.

**Exemple** :
```bash
dagger call \
  with-inventory --hosts "server1.example.com,server2.example.com" \
  with-inventory-var --key ansible_user --value ubuntu \
  with-inventory-var --key ansible_port --value 22 \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa
```

#### WithAnsibleVersion

Configure la version d'Ansible (défaut: `2.17`).

**Exemple** :
```bash
dagger call \
  with-ansible-version --version 2.16 \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa
```

#### WithTags

Configure les tags Ansible à exécuter.

**Exemple** :
```bash
dagger call \
  with-tags --tags "deploy,configure" \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa
```

#### WithSkipTags

Configure les tags Ansible à ignorer.

**Exemple** :
```bash
dagger call \
  with-skip-tags --tags "tests,slow" \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa
```

### Exécution

#### RunPlaybook

Exécute un playbook Ansible.

**Paramètres** :
- `playbook` : Chemin vers le playbook (ex: site.yml)
- `project` : Répertoire du projet Ansible
- `workdir` : Répertoire de travail relatif (défaut: `.`)
- `ssh-private-key` : Clé SSH privée (supporte `file:`, `env:`, etc.) **[OBLIGATOIRE]**
- `check-mode` : Mode dry-run (défaut: `false`)
- `verbose` : Niveau de verbosité 0-4 (défaut: `0`)
- `limit` : Limit pattern pour restreindre l'exécution

**Exemples** :
```bash
# Exécution simple
dagger call \
  with-inventory --hosts "server1.example.com" \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa

# Mode check (dry-run)
dagger call \
  with-inventory --hosts "server1.example.com" \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa --check-mode

# Avec verbosité
dagger call \
  with-inventory --hosts "server1.example.com" \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa --verbose 3

# Avec limit
dagger call \
  with-inventory --hosts "web1.example.com,web2.example.com,db1.example.com" \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa --limit "web*"
```

## 🔧 Exemples d'Utilisation

### Scénario 1 : Déploiement Simple

```bash
dagger call \
  with-ansible-version --version 2.17 \
  with-inventory --hosts "app1.example.com,app2.example.com" \
  with-inventory-var --key ansible_user --value deploy \
  with-inventory-var --key ansible_port --value 22 \
  with-variable --key app_version --value 1.2.3 \
  with-variable --key environment --value production \
  with-secret --key deploy_token --value env:DEPLOY_TOKEN \
  with-tags --tags "deploy,configure" \
  run-playbook --playbook deploy.yml --project /path/to/ansible --ssh-private-key file:~/.ssh/id_rsa
```

### Scénario 2 : Mode Check (Dry-Run)

```bash
dagger call \
  with-inventory --hosts "staging1.example.com,staging2.example.com" \
  with-variable --key ansible_user --value ansible \
  run-playbook \
    --playbook site.yml \
    --project . \
    --ssh-private-key file:~/.ssh/id_rsa \
    --check-mode \
    --verbose 2
```

### Scénario 3 : Pipeline CI/CD

```bash
# Étape 1 : Validation (check mode)
dagger call \
  with-inventory --hosts "prod1.example.com,prod2.example.com" \
  with-variable --key ansible_user --value deploy \
  run-playbook --playbook site.yml --project . --ssh-private-key env:SSH_PRIVATE_KEY --check-mode

# Étape 2 : Déploiement (si validation ok)
dagger call \
  with-inventory --hosts "prod1.example.com,prod2.example.com" \
  with-variable --key ansible_user --value deploy \
  with-extra-var --key deployment_id --value "$(date +%Y%m%d%H%M%S)" \
  with-tags --tags "deploy" \
  run-playbook --playbook site.yml --project . --ssh-private-key env:SSH_PRIVATE_KEY --verbose 1
```

## 🔐 Sécurité

### Gestion des Secrets

Les secrets ne sont **jamais** exposés dans les logs grâce à l'utilisation de `WithSecret` et des secrets Dagger :

```bash
# ✅ CORRECT : Secret protégé
dagger call \
  with-secret --key ansible_password --value env:MY_PASSWORD \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa

# ❌ INCORRECT : Valeur visible dans les logs
dagger call \
  with-variable --key ansible_password --value my-secret-password \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa
```

### SSH Keys

La clé SSH privée est obligatoire et doit être fournie via un secret Dagger :

```bash
# Depuis un fichier local
dagger call \
  run-playbook --playbook site.yml --project . --ssh-private-key file:~/.ssh/id_rsa

# Depuis une variable d'environnement
dagger call \
  run-playbook --playbook site.yml --project . --ssh-private-key env:SSH_PRIVATE_KEY
```

## 🧪 Tests

```bash
# Tester que le module est chargé
dagger call test

# Exécuter un playbook en mode check
dagger call \
  with-inventory --hosts "localhost" \
  run-playbook --playbook test.yml --project ./examples --ssh-private-key file:~/.ssh/id_rsa --check-mode
```

## 📚 Ressources

- [Documentation Ansible](https://docs.ansible.com)
- [Documentation Dagger](https://docs.dagger.io)
- [Dagger Secrets](https://docs.dagger.io/api/secrets)

## 🤝 Contribution

Pour ajouter une nouvelle fonctionnalité :

1. Créer un nouveau fichier (ex: `configure.go`)
2. Implémenter la fonction sur le type `*Ansible`
3. Utiliser les helpers partagés (`buildContainer`, `injectVariables`, etc.)
4. Ajouter la documentation dans ce README

Exemple :

```go
// configure.go
package main

import (
    "context"
    "dagger/ansible/internal/dagger"
)

func (m *Ansible) Configure(ctx context.Context, project *dagger.Directory, config string) (string, error) {
    container := m.buildContainer(project, ".")
    container, _ = m.injectVariables(ctx, container)
    return container.WithExec([]string{"ansible-config", config}).Stdout(ctx)
}
```

## 📄 Licence

MIT
