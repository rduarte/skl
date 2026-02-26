<p align="center">
  <h1 align="center">⚡ skl</h1>
  <p align="center">
    <strong>Gerenciador de skills de IA para projetos locais</strong>
  </p>
  <p align="center">
    Instale, atualize e gerencie capacidades de IA diretamente no seu projeto — como um <code>apt-get</code> para skills.
  </p>
</p>

---

## O que é o skl?

O **skl** é uma ferramenta de linha de comando que baixa e organiza _skills_ (capacidades e ferramentas de IA) dentro de projetos locais. Skills são armazenadas em repositórios Git (GitHub ou Bitbucket) e instaladas no diretório `.agent/skills/` do seu projeto.

```
seu-projeto/
├── .agent/
│   └── skills/
│       ├── 1doc-api-expert/    ← skill instalada
│       │   ├── SKILL.md
│       │   ├── docs/
│       │   └── knowledge/
│       └── data-analyzer/      ← outra skill
│           └── SKILL.md
├── sklfile.json                ← manifesto de dependências
└── ...
```

---

## Instalação

### Instalação rápida (recomendado)

```bash
curl -sSfL https://raw.githubusercontent.com/rduarte/skl/main/install.sh | bash
```

O script detecta a arquitetura, baixa o binário da última release e instala em `~/.local/bin/skl`.

> [!NOTE]
> Se `~/.local/bin` não estiver no seu `PATH`, adicione ao `~/.bashrc`:
> ```bash
> export PATH="$HOME/.local/bin:$PATH"
> ```

### Instalação manual

Baixe o binário diretamente na [página de releases](https://github.com/rduarte/skl/releases/latest):

```bash
# Baixar
wget https://github.com/rduarte/skl/releases/latest/download/skl-linux-amd64

# Instalar
chmod +x skl-linux-amd64
mv skl-linux-amd64 ~/.local/bin/skl
```

### Verificar instalação

```bash
skl --version
# skl version v0.1.0
```

---

## Atualização

O skl se atualiza sozinho:

```bash
skl upgrade
```

```
📦 Versão atual: v0.1.0
🔍 Verificando última versão...
⬆  Nova versão disponível: v0.2.0
⬇  Baixando skl-linux-amd64...
✅ skl atualizado para v0.2.0
```

---

## Uso

### Sintaxe geral

```
skl <comando> [argumentos]
```

A referência de uma skill segue o formato:

```
<provider>@<usuário>/<repositório>/<skill>[:tag]
```

| Componente     | Descrição                               | Exemplo             |
|----------------|-----------------------------------------|---------------------|
| `provider`     | Plataforma Git (`github`, `bitbucket`)  | `bitbucket`         |
| `usuário`      | Dono ou organização do repositório      | `servicos-1doc`     |
| `repositório`  | Nome do repositório                     | `1doc-apis`         |
| `skill`        | Nome da skill (subdiretório no repo)    | `1doc-api-expert`   |
| `tag`          | Versão específica _(opcional)_          | `v1.2.0`            |

---

## Comandos

### `skl install` — Instalar uma skill

Baixa uma skill de um repositório Git e a registra no manifesto do projeto.

```bash
# Instalar da branch padrão
skl install bitbucket@servicos-1doc/1doc-apis/1doc-api-expert

# Instalar uma versão específica
skl install github@empresa/repo-skills/data-analyzer:v1.2.0
```

**O que acontece:**
1. Clona o repositório via SSH (sparse-checkout — baixa **apenas** a skill)
2. Copia os arquivos para `.agent/skills/<skill>/`
3. Registra a dependência no `sklfile.json`

```
🔗 Clone URL: git@bitbucket.org:servicos-1doc/1doc-apis.git
⬇  Baixando skill "1doc-api-expert"...
✅ Skill "1doc-api-expert" instalada em .agent/skills/1doc-api-expert
📋 Skill registrada no sklfile.json
```

---

### `skl update` — Atualizar todas as skills

Lê o `sklfile.json` e instala ou atualiza todas as skills listadas.

```bash
skl update
```

```
📋 2 skill(s) encontrada(s) no sklfile.json

🔗 Clone URL: git@bitbucket.org:servicos-1doc/1doc-apis.git
⬇  Baixando skill "1doc-api-expert"...
✅ Skill "1doc-api-expert" instalada em .agent/skills/1doc-api-expert

🔗 Clone URL: git@github.com:empresa/repo-skills.git
⬇  Baixando skill "data-analyzer"...
✅ Skill "data-analyzer" instalada em .agent/skills/data-analyzer

📊 Resultado: 2/2 skill(s) instalada(s)
```

> [!TIP]
> Use `skl update` após clonar um projeto que tenha `sklfile.json` para instalar todas as skills de uma vez — semelhante a `npm install` ou `composer install`.

---

### `skl info` — Exibir informações de uma skill

Renderiza o `SKILL.md` de uma skill com formatação rica diretamente no terminal.

```bash
# Skill instalada localmente
skl info 1doc-api-expert

# Skill remota (sem instalar)
skl info bitbucket@servicos-1doc/1doc-apis/1doc-api-expert
skl info github@empresa/repo-skills/data-analyzer:v1.2.0
```

---

### `skl upgrade` — Atualizar o próprio skl

Verifica a última versão disponível no GitHub e atualiza o binário automaticamente.

```bash
skl upgrade
```

---

## Manifesto (`sklfile.json`)

O `sklfile.json` é o arquivo de manifesto que lista todas as skills do projeto. Ele é criado e atualizado automaticamente pelo comando `install`.

```json
{
  "skills": {
    "bitbucket@servicos-1doc/1doc-apis/1doc-api-expert": "*",
    "github@empresa/repo-skills/data-analyzer": "v1.2.0"
  }
}
```

| Valor     | Significado                                |
|-----------|--------------------------------------------|
| `"*"`     | Usa a branch padrão do repositório (latest)|
| `"v1.2.0"`| Versão fixa (tag Git)                     |

### Bloqueio de versões (`sklfile.lock`)

O `sklfile.lock` registra o estado exato das skills que estão instaladas. Ele é usado pelo comando `update` para calcular o diff entre o que você **deseja** (`sklfile.json`) e o que você **tem** (`sklfile.lock`).

**Por que o lock é importante?**
1. **Segurança**: Garante que todos os desenvolvedores do time tenham exatamente as mesmas versões.
2. **Sincronização**: Permite que o `skl update` remova automaticamente skills que foram deletadas do manifesto por outros desenvolvedores.

> [!IMPORTANT]
> Assim como no Composer (`composer.lock`) ou NPM (`package-lock.json`), você **deve** versionar o `sklfile.lock` no seu repositório.

---

## Pré-requisitos

- **Linux** (amd64)
- **Git** instalado e configurado com chave SSH para os repositórios desejados
- **Acesso SSH** aos repositórios que contêm as skills

> [!NOTE]
> O skl utiliza o protocolo SSH (`git@`) para clonagem, aproveitando as credenciais já configuradas no ambiente do usuário.

---

## Providers suportados

| Provider    | Clone URL                                    |
|-------------|----------------------------------------------|
| `github`    | `git@github.com:<user>/<repo>.git`           |
| `bitbucket` | `git@bitbucket.org:<user>/<repo>.git`        |

---

## Estrutura de uma skill no repositório

Para que o skl reconheça uma skill, ela deve estar localizada em:

```
<repositório>/
└── .agent/
    └── skills/
        └── <nome-da-skill>/
            ├── SKILL.md        ← obrigatório
            └── ...             ← outros arquivos da skill
```

---

## Referência rápida

```bash
# Instalar o skl
curl -sSfL https://raw.githubusercontent.com/rduarte/skl/main/install.sh | bash

# Instalar uma skill
skl install bitbucket@org/repo/skill-name

# Instalar todas as skills do projeto
skl update

# Ver informações de uma skill
skl info skill-name

# Atualizar o skl
skl upgrade

# Verificar versão
skl --version
```

---

## Licença

[CC0 1.0 Universal](LICENSE)
