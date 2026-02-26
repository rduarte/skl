# 📋 Notas de Lançamento (Release Notes)

Este arquivo documenta o histórico de versões do **skl**.

---

## [[v0.5.2]](https://github.com/rduarte/skl/releases/tag/v0.5.2) — 2026-02-26
### 🚀 Novas Funcionalidades
- **Suporte a Múltiplos Caminhos de Skill**: O `skl` agora detecta automaticamente skills localizadas tanto em `.agent/skills/` quanto na pasta `skills/` (na raiz do repositório remoto).
- Facilita a integração com repositórios que seguem diferentes convenções de organização.

---

## [[v0.5.1]](https://github.com/rduarte/skl/releases/tag/v0.5.1) — 2026-02-26
### 📝 Documentação
- Adicionado link para o `RELEASE_NOTES.md` no final do `README.md` para facilitar o acesso ao histórico de mudanças.

---

## [[v0.5.0]](https://github.com/rduarte/skl/releases/tag/v0.5.0) — 2026-02-26
### 🚀 Novas Funcionalidades
- **Detecção Inteligente de Atualizações (Smart Update)**: O `skl` agora rastreia o hash exato do commit no `sklfile.lock`.
- O comando `update` agora detecta mudanças no repositório remoto mesmo usando referências simbólicas como `*` ou nomes de branch.
- Resolução dinâmica de hashes via `git ls-remote` (sem necessidade de download prévio).

---

## [[v0.4.11]](https://github.com/rduarte/skl/releases/tag/v0.4.11) — 2026-02-26
### 📝 Documentação
- Migração das notas de lançamento para o arquivo local `RELEASE_NOTES.md`.
- Adição de links diretos para as tags do GitHub no histórico.

---

## [[v0.4.10]](https://github.com/rduarte/skl/releases/tag/v0.4.10) — 2026-02-26
### 🐛 Correções de Bugs
- Corrigida sintaxe no esquema de configuração de releases do GitHub (removido em favor deste arquivo).

---

## [[v0.4.9]](https://github.com/rduarte/skl/releases/tag/v0.4.9) — 2026-02-26
### ⚙️ Manutenção
- Adição inicial de configuração de release do GitHub (substituída por este arquivo local).

---

## [[v0.4.8]](https://github.com/rduarte/skl/releases/tag/v0.4.8) — 2026-02-26
### 📝 Documentação
- Reorganização lógica do README.md focada em cenários de uso (Novo Projeto, Adoção Local e Time Sync).

---

## [[v0.4.7]](https://github.com/rduarte/skl/releases/tag/v0.4.7) — 2026-02-26
### 🐛 Correções de Bugs
- Tratamento amigável para erro 403 (Rate Limit) da API do GitHub no comando `upgrade`.

---

## [[v0.4.6]](https://github.com/rduarte/skl/releases/tag/v0.4.6) — 2026-02-26
### 🚀 Novas Funcionalidades
- **Novo Comando `setup`**: Indexação automática de pastas locais em `.agent/skills` no manifesto.
- **Suporte a `local@`**: Novo prefixo para gerenciar skills sem repositório remoto.
### ⚙️ Manutenção e Refatoração
- Renomeado antigo comando `setup` para `rebuild` (focado apenas em autocomplete).

---

## [[v0.4.5]](https://github.com/rduarte/skl/releases/tag/v0.4.5) — 2026-02-26
### 📝 Documentação
- Overhaul completo do README.md e CONTRIBUTING.md para maior clareza e segmentação.

---

## [[v0.4.4]](https://github.com/rduarte/skl/releases/tag/v0.4.4) — 2026-02-26
### 🚀 Novas Funcionalidades
- Autocomplete dinâmico para os comandos `info` e `remove` (baseado no `sklfile.lock`).

---

## [[v0.4.3]](https://github.com/rduarte/skl/releases/tag/v0.4.3) — 2026-02-26
### 🐛 Correções de Bugs
- Abortar `skl update` com erro claro se o manifesto `sklfile.json` estiver ausente.

---

## [[v0.4.2]](https://github.com/rduarte/skl/releases/tag/v0.4.2) — 2026-02-26
### 🚀 Novas Funcionalidades
- Autocomplete agora funciona sem a necessidade de barra final (`/`) no repositório.

---

## [[v0.4.1]](https://github.com/rduarte/skl/releases/tag/v0.4.1) — 2026-02-26
### ⚙️ Manutenção e Refatoração
- Configuração de autocomplete agora é totalmente silenciosa e shell-aware.

---

## [[v0.4.0]](https://github.com/rduarte/skl/releases/tag/v0.4.0) — 2026-02-26
### 🚀 Novas Funcionalidades
- Suporte ao campo `path` no `catalog.json`, permitindo redirecionar para a localização exata da skill no repositório.
- Adição da flag `--force` no comando `install`.

---

## [[v0.3.3]](https://github.com/rduarte/skl/releases/tag/v0.3.3) — 2026-02-25
### 🚀 Novas Funcionalidades
- Verificação automática de novas versões do skl em cada execução.

---

## [[v0.3.2]](https://github.com/rduarte/skl/releases/tag/v0.3.2) — 2026-02-25
### 🚀 Novas Funcionalidades
- Implementação do comando `skl list` para explorar repositórios.

---

## [[v0.3.1]](https://github.com/rduarte/skl/releases/tag/v0.3.1) — 2026-02-25
### 🚀 Novas Funcionalidades
- Implementação de autocompletar nativo para Bash e Zsh.

---

## [[v0.3.0]](https://github.com/rduarte/skl/releases/tag/v0.3.0) — 2026-02-25
### 🚀 Novas Funcionalidades
- Suporte inicial a catálogos (`catalog.json`) para descoberta de skills.

---

## [[v0.2.0]](https://github.com/rduarte/skl/releases/tag/v0.2.0) — 2026-02-25
### 🚀 Novas Funcionalidades
- **Sincronização Determinística**: Implementação do `sklfile.lock`.
- Lógica de atualização baseada em diff entre manifesto e estado atual.

---

## [[v0.1.0]](https://github.com/rduarte/skl/releases/tag/v0.1.0) — 2026-02-25
### 🚀 Novas Funcionalidades
- Lançamento inicial do **skl**.
- Comandos base: `install`, `update`, `info`, `upgrade`, `setup`.
