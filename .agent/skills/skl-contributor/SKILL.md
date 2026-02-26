---
name: SKLContributor
description: Expert on contributing to the skl project — architecture, development workflow, adding providers, and publishing releases.
---

# Guia do Contribuidor skl

Este documento é a referência para contribuir com o projeto **skl** — o gerenciador de skills de IA.

> [!IMPORTANT]
> Para detalhes técnicos completos sobre configuração de ambiente, criação de providers e workflow detalhado, consulte o arquivo [CONTRIBUTING.md](file:///home/ricardo/git/skl/CONTRIBUTING.md) na raiz do projeto.

## Visão Geral da Arquitetura

O **skl** é escrito em Go e segue uma estrutura modular:

```
main.go → cmd/root.go → cmd/{install,update,info,upgrade}.go
                              │
                              ▼
                    internal/parser/parser.go      ← parsing de referências
                    internal/provider/provider.go  ← interface + registry
                    internal/installer/installer.go ← sparse-checkout + copy
                    internal/manifest/manifest.go   ← sklfile.json & lock
```

### Fluxo de Dados

1. **Parsing**: `<provider>@<user>/<repo>/<skill>[:tag]` → struct `SkillRef`
2. **Provider**: Resolve URLs de clone e download.
3. **Installer**: Realiza `sparse-checkout` e organiza os arquivos em `.agent/skills/`.
4. **Manifest / Lock**: Garante a persistência e o versionamento determinístico (Smart Update).

## Fluxo de Release

As releases são automatizadas via GitHub Actions.

1. **Quando Lançar?**: Crie uma nova versão apenas para alterações funcionais (novos comandos, bugs, mudanças no binário).
   - **Documentação**: Alterações em `README.md` ou `CONTRIBUTING.md` devem ser commitadas no `main` **sem** gerar uma nova tag.
2. **Tagging**:
   ```bash
   git tag v0.5.1
   git push origin v0.5.1
   ```
3. **Automação**: O workflow de CI compila o binário e o comando `skl upgrade` dos usuários detectará a nova versão automaticamente.

## Guia de Desenvolvimento

Consulte o [CONTRIBUTING.md](file:///home/ricardo/git/skl/CONTRIBUTING.md) para:
- Como adicionar novos **Providers**.
- Como adicionar novos **Comandos**.
- Padrões de **Commits Semânticos**.
- Comandos de `Makefile` e automação de build.

## Histórico de Mudanças

Consulte o [RELEASE_NOTES.md](file:///home/ricardo/git/skl/RELEASE_NOTES.md) para o registro completo de todas as versões lançadas.
