---
name: SKLContributor
description: Expert on contributing to the skl project — architecture, development workflow, adding providers, and publishing releases.
---

# Guia do Contribuidor skl

Este documento é a referência para contribuir com o projeto **skl** — o gerenciador de skills de IA.

## Visão Geral da Arquitetura

```
main.go → cmd/root.go → cmd/{install,update,info,upgrade}.go
                              │
                              ▼
                    internal/parser/parser.go      ← parsing de referências
                    internal/provider/provider.go  ← interface + registry
                    internal/provider/{github,bitbucket}.go
                    internal/installer/installer.go ← sparse-checkout + copy
                    internal/manifest/manifest.go   ← sklfile.json
```

### Fluxo de Dados

1. **Parsing**: `<provider>@<user>/<repo>/<skill>[:tag]` → struct `SkillRef`
2. **Provider**: `SkillRef.Provider` → `CloneURL()` e `RepoURL()`
3. **Installer**: sparse-checkout do repo → copia `.agent/skills/<skill>/` para destino
4. **Manifest**: registra/lê skills no `sklfile.json`

## Decisões de Design

| Decisão | Regra |
|---|---|
| Skills no repo | Sempre em `.agent/skills/<nome>/` |
| Protocolo Git | SSH (`git@`) — nunca HTTPS |
| Sparse-checkout | `--filter=blob:none --sparse --depth=1 --no-checkout` |
| Validação | `git ls-tree HEAD <path>` antes do checkout |
| Manifesto | `sklfile.json` na raiz — chave = referência completa, valor = tag ou `*` |
| Versão | Injetada via `-ldflags -X cmd.Version` no build |
| Erros do git | Stderr capturado e classificado (repo não encontrado, acesso negado, tag inválida) |

## Como Adicionar um Provider

1. Criar `internal/provider/<nome>.go` implementando a interface:

```go
type Provider interface {
    Name() string
    CloneURL(user, repo string) string
    RepoURL(user, repo string) string
}
```

2. Registrar no mapa `registry` em `internal/provider/provider.go`
3. Nenhuma outra mudança necessária — o factory `New()` resolve automaticamente

## Como Adicionar um Comando

1. Criar `cmd/<comando>.go`
2. Definir `var <comando>Cmd = &cobra.Command{...}`
3. Registrar no `init()` com `rootCmd.AddCommand(<comando>Cmd)`

## Fluxo de Release

```bash
# 1. Commitar mudanças (Conventional Commits)
git commit -m "feat: nova funcionalidade"

# 2. Criar tag semântica
git tag v0.2.0

# 3. Push (dispara GitHub Actions)
git push origin main --tags
```

O workflow `.github/workflows/release.yml`:
- Compila `CGO_ENABLED=0 GOOS=linux GOARCH=amd64`
- Injeta versão via ldflags
- Cria Release com binário `skl-linux-amd64`

## Convenção de Commits

| Prefixo | Uso |
|---|---|
| `feat:` | Nova funcionalidade |
| `fix:` | Correção de bug |
| `docs:` | Documentação |
| `refactor:` | Refatoração sem mudar comportamento |
| `chore:` | Manutenção e build |

## Comandos de Desenvolvimento

```bash
make build                    # Compilar binário
make build VERSION=v0.2.0     # Compilar com versão específica
make install                  # Instalar em ~/.local/bin
make clean                    # Remover binário
```
