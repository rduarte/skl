# Contribuindo com o skl

Obrigado pelo interesse em contribuir! Este guia explica como configurar o ambiente, desenvolver, e publicar novas versões do skl.

---

## Pré-requisitos

- **Go 1.24+** — [go.dev/dl](https://go.dev/dl/)
- **Git** com chave SSH configurada
- **Make** (geralmente já instalado no Linux)

---

## Configuração do ambiente

```bash
# 1. Clone o repositório
git clone git@github.com:rduarte/skl.git
cd skl

# 2. Instale as dependências
go mod download

# 3. Compile
make build

# 4. Verifique
./skl --version
```

---

## Estrutura do projeto

```
skl/
├── main.go                          # Ponto de entrada
├── cmd/
│   ├── root.go                      # Comando raiz (Cobra)
│   ├── install.go                   # skl install
│   ├── update.go                    # skl update
│   ├── info.go                      # skl info
│   └── upgrade.go                   # skl upgrade (self-update)
├── internal/
│   ├── parser/parser.go             # Parser de referências de skills
│   ├── provider/
│   │   ├── provider.go              # Interface + registry de providers
│   │   ├── github.go                # Provider GitHub
│   │   └── bitbucket.go             # Provider Bitbucket
│   ├── installer/installer.go       # Lógica de clone + sparse-checkout
│   └── manifest/manifest.go         # Leitura/gravação do sklfile.json
├── install.sh                       # Script de instalação para usuários
├── Makefile                         # Build, install, clean
├── .github/workflows/release.yml    # CI/CD para releases automáticas
└── go.mod / go.sum
```

---

## Fluxo de desenvolvimento

### 1. Crie uma branch

```bash
git checkout -b feat/minha-feature
```

### 2. Faça suas alterações

- Novos comandos vão em `cmd/`
- Lógica interna vai em `internal/`
- Novos providers implementam a interface `Provider` em `internal/provider/`

### 3. Compile e teste

```bash
# Build rápido
make build
./skl --help

# Testar um comando
./skl install github@user/repo/skill
```

### 4. Commit e push

Use [Conventional Commits](https://www.conventionalcommits.org/):

```bash
git commit -m "feat: adicionar suporte a GitLab"
git commit -m "fix: tratar erro de timeout no clone"
git commit -m "docs: atualizar exemplos no README"
```

| Prefixo  | Quando usar                        |
|----------|------------------------------------|
| `feat:`  | Nova funcionalidade                |
| `fix:`   | Correção de bug                    |
| `docs:`  | Alteração em documentação          |
| `refactor:` | Refatoração sem mudar comportamento |
| `chore:` | Tarefas de manutenção              |

### 5. Abra um Pull Request

```bash
git push origin feat/minha-feature
```

---

## Adicionando um novo provider

Para suportar uma nova plataforma Git (ex: GitLab):

1. Crie `internal/provider/gitlab.go`:

```go
package provider

import "fmt"

type GitLab struct{}

func (GitLab) Name() string { return "gitlab" }

func (GitLab) CloneURL(user, repo string) string {
    return fmt.Sprintf("git@gitlab.com:%s/%s.git", user, repo)
}

func (GitLab) RepoURL(user, repo string) string {
    return fmt.Sprintf("https://gitlab.com/%s/%s", user, repo)
}
```

2. Registre no `provider.go`:

```go
var registry = map[string]Provider{
    "github":    GitHub{},
    "bitbucket": Bitbucket{},
    "gitlab":    GitLab{},   // ← adicionar aqui
}
```

3. Compile e teste:

```bash
make build
./skl install gitlab@grupo/repo/skill
```

---

## Publicando uma nova release

### 1. Atualize a versão

Defina a tag semântica seguindo [SemVer](https://semver.org/):

| Tipo de mudança             | Exemplo       |
|-----------------------------|---------------|
| Correção de bug             | `v0.1.1`      |
| Nova funcionalidade         | `v0.2.0`      |
| Breaking change             | `v1.0.0`      |

### 2. Crie a tag e faça push

```bash
git tag v0.2.0
git push origin main --tags
```

### 3. Release automática

O **GitHub Actions** faz o resto automaticamente:
1. Compila o binário para `linux/amd64` com a versão embutida
2. Cria a GitHub Release com o binário anexado
3. Release notes são geradas automaticamente

### 4. Verifique

Acesse: https://github.com/rduarte/skl/releases

Usuários finais atualizam com:
```bash
skl upgrade
```

---

## Build com versão customizada

```bash
# Build local com versão
make build VERSION=v0.2.0-beta

# Verificar
./skl --version
# skl version v0.2.0-beta
```

A versão é injetada via `-ldflags` no build:
```
-X github.com/rduarte/skl/cmd.Version=$(VERSION)
```

---

## Arquitetura de decisões

| Decisão | Justificativa |
|---|---|
| **Cobra** para CLI | Framework padrão do ecossistema Go, auto gera help e completions |
| **Sparse-checkout** | Baixa apenas o diretório da skill, não o repo inteiro |
| `--depth=1` + `--filter=blob:none` | Minimiza tráfego de rede e uso de disco |
| **SSH por padrão** | Aproveita credenciais já configuradas no ambiente do dev |
| **Provider registry** | Adicionar provider = 1 arquivo, zero mudanças no código existente |
| **`sklfile.json`** | Manifesto simples e legível, inspirado em composer.json |
| **Self-update** | Consulta GitHub API, baixa binário e substitui in-place |
