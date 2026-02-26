<p align="center">
  <h1 align="center">⚡ skl</h1>
  <p align="center">
    <strong>O gerenciador de dependências para Skills de IA.</strong>
  </p>
  <p align="center">
    Instale, organize e compartilhe capacidades de IA diretamente no seu projeto.
  </p>
</p>

---

## 🚀 O que é o skl?

O **skl** é um gerenciador de pacotes especializado em _skills_ (conhecimento, ferramentas e personalidades de IA). Ele permite que você baixe scripts, configurações e documentação de repositórios Git externos e os organize automaticamente no seu workflow local.

As skills são instaladas no diretório `.agent/skills/` e o projeto mantém um manifesto (`sklfile.json`) para que qualquer membro do time possa sincronizar o ambiente com um único comando.

---

## 🛠️ Instalação

### Instalação Rápida (Linux / macOS)

```bash
curl -sSfL https://raw.githubusercontent.com/rduarte/skl/main/install.sh | bash
```

O script instala o binário em `~/.local/bin/skl` e configura o autocompletar automaticamente para seu shell (Bash ou Zsh).

> [!TIP]
> **Autocomplete:** Após instalar, abra um novo terminal para ativar as sugestões de comando via `TAB`.

---

## 📖 Guia de Uso

### 1. Explorando Repositórios
Você pode listar as skills disponíveis em um repositório antes de instalar:

```bash
skl list github@rmyndharis/antigravity-skills
```

### 2. Instalando Skills
Para instalar uma skill, use o formato `provedor@usuario/repo/skill`:

```bash
# Instalando a versão mais recente
skl install github@rmyndharis/antigravity-skills/tutorial-engineer

# Forçando a reinstalação de uma skill existente
skl install github@user/repo/skill --force
```

### 3. Sincronizando o Projeto
Ao entrar em um projeto que já possui um `sklfile.json`, basta rodar:

```bash
skl update
```
Isso instalará todas as skills faltantes e removerá as que não estão mais no manifesto.

### 4. Consultando Documentação
Cada skill possui um arquivo `SKILL.md`. Você pode lê-lo formatado no terminal:

```bash
# Skill instalada
skl info tutorial-engineer

# Skill remota (sem precisar instalar)
skl info github@rmyndharis/antigravity-skills/ai-expert
```

### 5. Removendo Skills
```bash
skl remove tutorial-engineer
```

---

## 📋 Arquivos de Configuração

### `sklfile.json` (Manifesto)
Lista o que seu projeto **deseja** ter. Deve ser compartilhado com seu time.

### `sklfile.lock` (Estado Atual)
Registra o que está **realmente instalado** (versões exatas, tags, etc). Garante que todos os desenvolvedores usem a mesma versão das ferramentas.

---

## ⚙️ Comandos Disponíveis

| Comando | Descrição |
| :--- | :--- |
| `install` | Baixa e registra uma nova skill. |
| `update` | Sincroniza o projeto local com o manifesto. |
| `remove` | Exclui uma skill local e a remove do manifesto. |
| `list` | Lista skills disponíveis em um repositório remoto. |
| `info` | Exibe a documentação (`SKILL.md`) da skill. |
| `upgrade` | Atualiza o próprio `skl` para a última versão. |
| `rebuild` | (Automático) Reconstrói as sugestões de shell. |
| `setup` | Indexa folders locais como skills gerenciadas. |

---

## 🤝 Contribuindo

Interessado em adicionar novos providers, comandos ou melhorar a arquitetura? Veja nosso [Guia de Contribuição](CONTRIBUTING.md).

---

## 📄 Licença

[CC0 1.0 Universal](LICENSE)
