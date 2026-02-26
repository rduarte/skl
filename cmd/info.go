package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/rduarte/skl/internal/installer"
	"github.com/rduarte/skl/internal/parser"
	"github.com/rduarte/skl/internal/provider"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <skill-name | provider@user/repo/skill[:tag]>",
	Short: "Exibe o SKILL.md de uma skill instalada ou remota",
	Long: `Lê e renderiza o arquivo SKILL.md formatado para o terminal.

Uso com skill instalada:
  skl info 1doc-api-expert

Uso com referência remota (sem instalar):
  skl info bitbucket@servicos-1doc/1doc-apis/1doc-api-expert
  skl info github@empresa/repo-skills/data-analyzer:v1.2.0`,
	Args: cobra.ExactArgs(1),
	RunE: runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) error {
	arg := args[0]

	var data []byte
	var err error

	if strings.Contains(arg, "@") {
		// Remote mode: fetch SKILL.md from repo without installing
		data, err = fetchRemoteSkillMD(arg)
	} else {
		// Local mode: read from installed skill
		data, err = readLocalSkillMD(arg)
	}

	if err != nil {
		return err
	}

	return renderMarkdown(data)
}

// readLocalSkillMD reads SKILL.md from a locally installed skill.
func readLocalSkillMD(skill string) ([]byte, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter diretório atual: %w", err)
	}

	skillFile := filepath.Join(cwd, ".agent", "skills", skill, "SKILL.md")

	data, err := os.ReadFile(skillFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("skill %q não encontrada\n\n  Arquivo esperado: %s\n  A skill está instalada? Verifique com: skl update", skill, skillFile)
		}
		return nil, fmt.Errorf("erro ao ler SKILL.md: %w", err)
	}

	return data, nil
}

// fetchRemoteSkillMD fetches SKILL.md from a remote repo using sparse-checkout.
func fetchRemoteSkillMD(rawRef string) ([]byte, error) {
	ref, err := parser.Parse(rawRef)
	if err != nil {
		return nil, err
	}

	prov, err := provider.New(ref.Provider)
	if err != nil {
		return nil, err
	}

	cloneURL := prov.CloneURL(ref.User, ref.Repo)
	repoURL := prov.RepoURL(ref.User, ref.Repo)

	fmt.Printf("🔗 Clone URL: %s\n", cloneURL)
	fmt.Printf("⬇  Buscando SKILL.md de %q...\n\n", ref.Skill)

	data, err := installer.FetchFile(cloneURL, repoURL, ref.Skill, ref.Tag, "SKILL.md")
	if err != nil {
		return nil, err
	}

	return data, nil
}

// renderMarkdown renders markdown content to the terminal.
func renderMarkdown(data []byte) error {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		fmt.Println(string(data))
		return nil
	}

	rendered, err := renderer.Render(string(data))
	if err != nil {
		fmt.Println(string(data))
		return nil
	}

	fmt.Print(rendered)
	return nil
}
