# Instruções do Agente para o Projeto `skl`

Você é um expert contribuindo para o projeto **skl** — o gerenciador de skills de IA.
Seu objetivo é auxiliar no desenvolvimento, manutenção e evolução desta ferramenta, seguindo rigorosamente as diretrizes arquiteturais e de fluxo de trabalho do projeto.

## 🏗️ Arquitetura do Projeto

O **skl** é escrito em Go e possui a seguinte estrutura modular:
- `cmd/`: Comandos CLI usando Cobra (`root.go`, `install.go`, `update.go`, `info.go`, `list.go`, etc.).
- `internal/parser/`: Lógica de parsing de referências (`provider@user/repo/skill:tag`) e repositórios.
- `internal/provider/`: Abstração de Git Hosts (GitHub, Bitbucket, Local). Interface principal para URLs de clone e raw.
- `internal/catalog/`: Busca e parse do `catalog.json` remoto via HTTP.
- `internal/installer/`: Lógica de `sparse-checkout` (via git) e cópia/gestão de arquivos em `.agent/skills/`.
- `internal/manifest/`: Gestão de estado através do `sklfile.json` (declarativo) e `sklfile.lock` (resolução exata de hashes, Smart Update).
- `internal/updater/`: Lógica de auto-update verificando GitHub Releases.

Fluxo de execução típico de um comando (ex: `install`):
`Parsing da URL` -> `Provider (Resolve URLs)` -> `Busca do Catálogo (Opcional)` -> `Installer (sparse-checkout)` -> `Manifest/Lock (Atualiza estado local)`.

## 🛠️ Regras de Desenvolvimento

1. **Testes Necessários**: Sempre execute testes antes de confirmar qualquer alteração.
2. **Commits Semânticos e Atômicos**: Use prefixos claros: `feat:`, `fix:`, `refactor:`, `docs:`, `chore:`. Cada commit deve representar uma única alteração lógica (atômico).
3. **Timeouts**: Comandos que fazem requisições de rede (como HTTP requests em autocompletes ou listagens) devem sempre respeitar timeouts curtos (ex: 2s a 10s) para não travar a experiência do usuário.
4. **Respostas Silenciosas**: Comandos de setup ou utilitários atuando em background (como autocompletar) não devem emitir logs desnecessários.
5. **Comandos Essenciais**:
   - Compilar o projeto: `go build -o /tmp/skl .` (útil para testar que o código compila).
   - Testar o binário: `go run main.go <comando> [args]`
6. **Autogerenciamento de Conhecimento**: Todo conhecimento relevante, decisões arquiteturais ou padrões identificados no decorrer das conversas ao atuar neste projeto devem ser **automaticamente incorporados neste arquivo (`.agent/agent.md`)**. O agente deve ser proativo em manter sua própria documentação de base de conhecimento atualizada.

## 📦 Fluxo de Releases (Workflow `nova-versao`)

Temos um workflow padronizado para lançar novas versões. SEMPRE que o usuário solicitar uma nova versão, consulte e execute rigorosamente o workflow localizado em `.agent/workflows/nova-versao.md`.

*Nota de Histórico: Estas diretrizes absorveram e substituem o antigo skill `skl-contributor`.*
