# Guia de Contribuição ✨

Este guia é destinado a desenvolvedores que desejam evoluir o **skl**, adicionar novos providers ou entender as entranhas da ferramenta.

---

## 🏗️ Arquitetura do Projeto

O projeto é escrito em Go e segue uma estrutura modular para facilitar a expansão:

```
skl/
├── cmd/                # Comandos CLI (Cobra)
│   ├── root.go         # Configuração base + Version Check
│   ├── install.go      # Resolução via Catalog + Installer
│   ├── rebuild.go      # (Antigo setup) Reconstrói completions
│   ├── setup.go        # (Novo) Indexa skills locais
│   └── ...             # Demais comandos (update, remove, list, info)
├── internal/
│   ├── parser/         # Lógica de parsing de referências e repositórios
│   ├── provider/       # Abstração de Git Hosts (GitHub, Bitbucket)
│   ├── catalog/        # Busca e parse de catalog.json via HTTP
│   ├── installer/      # Clone, sparse-checkout e gestão de arquivos
│   ├── manifest/       # Gestão do sklfile.json e sklfile.lock
│   └── updater/        # Lógica de auto-update (GitHub Releases)
├── install.sh          # Script de instalação para usuário final
└── .github/workflows/  # CI/CD (Build e Release automática)
```

---

## 🛠️ Configuração do Ambiente

1. **Go 1.24+**: Certifique-se de ter o Go instalado.
2. **Clone**: `git clone git@github.com:rduarte/skl.git`
3. **Build para Teste**:
   ```bash
   go build -o skl .
   ./skl list github@rmyndharis/antigravity-skills
   ```

---

## 🚀 Como Criar um Novo Provider

Se você deseja adicionar suporte a uma nova plataforma (ex: GitLab), siga estes passos:

1. **Implemente a interface `Provider`** em `internal/provider/`:
   ```go
   type Provider interface {
       Name() string
       CloneURL(user, repo string) string
       RepoURL(user, repo string) string
       RawURL(user, repo, ref, path string) string // Para busca de catalog.json
   }
   ```
2. **Registre o provider** no mapa `registry` em `internal/provider/provider.go`.

---

## 📦 Fluxo de Release

As releases são automatizadas via GitHub Actions.

1. **Quando Lançar?**: Lance uma nova versão apenas para alterações que impactem o **uso ou comportamento** da ferramenta CLI (novos comandos, correções de bugs, features).
   - **Exemplo**: Se você apenas atualizou documentos (`README`, `CONTRIBUTING`), envie o commit para o `main`, mas **não crie uma nova tag**.
2. **Tagging**: Crie uma tag seguindo o versionamento semântico:
   ```bash
   git tag v0.4.5
   git push origin v0.4.5
   ```
3. **Automação**: O workflow de CI irá compilar o binário e criar a release no GitHub. O comando `skl upgrade` dos usuários detectará a nova versão automaticamente.

---

## 🧪 Boas Práticas

- **Commits Semânticos**: Use `feat:`, `fix:`, `refactor:`, `docs:` para manter o histórico organizado.
- **Timeouts**: Comandos que fazem rede (como `list` ou `upgrade`) devem sempre respeitar os timeouts definidos (geralmente entre 1.5s e 10s) para não travar a experiência do usuário.
- **Silêncio é Ouro**: Comandos de automação (como `setup` ou `upgrade`) devem ser o mais silenciosos possível, imprimindo apenas o estritamente necessário.

---

## 💡 Sugestões de Melhorias?

Abra uma **Issue** ou um **Pull Request**. Valorizamos simplicidade, velocidade e design minimalista. ⚡
