---
description: Lançar Nova Versão do skl
---

# Lançar Nova Versão do skl

Este workflow deve ser seguido RIGOROSAMENTE sempre que o usuário solicitar o lançamento de uma nova versão do projeto `skl`.

## Regras Absolutas

- Sempre a **documentação deve ser revisada e atualizada** antes de fechar a versão (ex: `README.md`, `RELEASE_NOTES.md`).
- Sempre os **testes devem ser executados** (ex: `go build` ou `go test`, dependendo do projeto) garantindo que o ramo atual está estável.
- Sempre os **commits devem ser feitos de maneira atômica e semântica** (ex: `feat:`, `fix:`).
- Sempre a **versão deve ser incrementada automaticamente (como PATCH)**, a menos que o usuário oriente expressamente algo diferente (ex: Minor ou Major).
- Sempre **deve ser criada uma tag no git** e enviada para o repositório junto com todos os commits. O CI/CD do GitHub Actions se encarregará da automação real da release.

---

## Passos da Execução

1. **Rever Estado e Documentação**:
   - Verifique o que mudou desde a última release (ex: via `git log` ou estado atual das modificações não commitadas).
   - Atualize impreterivelmente o topo do arquivo `RELEASE_NOTES.md` com as novas contribuições, detalhando Bugs ou Features.

2. **Garantir a Estabilidade (Testes)**:
   - Verifique que o projeto compila e passa nas validações básicas.
   ```bash
   go build -o /dev/null main.go
   go test ./...
   ```

3. **Definir e Incrementar a Versão**:
   - Descubra a tag mais recente:
   ```bash
   git describe --tags --abbrev=0
   ```
   - Incremente para a próxima versão de **Patch** (ex: `v0.5.5` -> `v0.5.6`), a menos que instruído em contrário pelo usuário. O título da seção no `RELEASE_NOTES.md` deve corresponder a essa nova versão.

4. **Commits Atômicos e Semânticos**:
   - Certifique-se de que os arquivos foram commitados com base no que representam.
   ```bash
   git add RELEASE_NOTES.md (e outros arquivos se existirem)
   git commit -m "chore(release): bump version to vX.Y.Z"
   ```

5. **Criar a Tag e Enviar para o Remoto**:
   - Crie a tag associada à versão e submeta tudo para origin.
   ```bash
   git tag vX.Y.Z
   git push origin main
   git push origin vX.Y.Z
   ```

6. **Notificar o Usuário**:
   - Avise o usuário de que a versão foi devidamente comitada, a Tag foi criada e enviada ao servidor, e que o pipeline de Release cuidará da materialização dos binários no GitHub.
