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

---

## 🛠️ Instalação

```bash
curl -sSfL https://raw.githubusercontent.com/rduarte/skl/main/install.sh | bash
```

O script instala o binário em `~/.local/bin/skl` e configura o autocompletar automaticamente para seu shell (Bash ou Zsh).

---

## 📖 Guia de Uso

Dependendo do estado do seu projeto, o fluxo de trabalho pode variar:

### A. Novo Projeto (Começando do zero)
Se você quer adicionar novas capacidades ao seu projeto:

1. **Explore as skills** disponíveis em um repositório:
   ```bash
   skl list github@rmyndharis/antigravity-skills
   ```
2. **Instale a skill** desejada:
   ```bash
   skl install github@rmyndharis/antigravity-skills/tutorial-engineer
   ```
   *Isso criará a pasta `.agent/skills/` e o arquivo `sklfile.json`.*

### B. Projeto com Skills Existentes (Adoção)
Se você já possui pastas de skills dentro de `.agent/skills/` (criadas manualmente ou legadas) e quer que o `skl` passe a gerenciá-las:

1. **Indexe as pastas locais**:
   ```bash
   skl setup
   ```
   *O `skl` detectará as pastas e as adicionará ao manifesto como `local@nome-da-skill`.*
2. **Pronto!** Agora o `skl` sabe que essas skills existem e não as removerá durante sincronizações.

### C. Trabalhando em Time (Manifesto existente)
Se você acabou de clonar um projeto que já possui um `sklfile.json`:

1. **Sincronize o ambiente**:
   ```bash
   skl update
   ```
   *Isso baixará todas as skills listadas e removerá qualquer uma que tenha sido deletada do manifesto.*

---

## ⚙️ Comandos Essenciais

| Comando | Descrição |
| :--- | :--- |
| `list` | Lista skills disponíveis em um repositório remoto. |
| `install` | Baixa e registra uma nova skill no projeto. |
| `setup` | Indexa diretórios locais em `.agent/skills` no manifesto. |
| `update` | Sincroniza as skills locais com o manifesto (`sklfile.json`). |
| `info` | Exibe a documentação (`SKILL.md`) da skill (local ou remota). |
| `remove` | Exclui uma skill e a remove do manifesto. |
| `upgrade` | Atualiza o próprio `skl` para a última versão. |

---

## 📋 Arquivos de Configuração

- **`sklfile.json`**: O manifesto de dependências. Lista o que seu projeto "deseja" ter.
- **`sklfile.lock`**: O registro do estado atual. Garante que todos no time usem as mesmas versões exatas.

---

## 🤝 Contribuindo

Para desenvolvedores que desejam evoluir a ferramenta ou adicionar novos providers, veja o [CONTRIBUTING.md](CONTRIBUTING.md).

---

## 📅 Histórico de Versões

Para acompanhar as últimas melhorias e correções, consulte nosso [RELEASES.md](RELEASES.md).

---

## 📄 Licença

[CC0 1.0 Universal](LICENSE)
