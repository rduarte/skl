package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/rduarte/skl/internal/catalog"
	"github.com/rduarte/skl/internal/installer"
	"github.com/rduarte/skl/internal/manifest"
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
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		var suggestions []string

		if strings.Contains(toComplete, "@") {
			parts := strings.Split(toComplete, "/")

			var ref *parser.SkillRef
			var refStr string

			if len(parts) >= 3 {
				// case: provider@user/repo/skill-prefix
				refStr = strings.Join(parts[:len(parts)-1], "/") + "/"
				r, err := parser.Parse(refStr + "dummy")
				if err == nil {
					ref = r
				}
			} else if len(parts) == 2 {
				// case: provider@user/repo (missing trailing slash)
				// we try to parse it as a repo to get user/repo
				r, err := parser.ParseRepo(toComplete)
				if err == nil {
					refStr = toComplete + "/"
					ref = &parser.SkillRef{
						Provider: r.Provider,
						User:     r.User,
						Repo:     r.Repo,
						Tag:      r.Tag,
					}
				}
			}

			if ref != nil {
				prov, err := provider.New(ref.Provider)
				if err == nil {
					cat, err := catalog.Fetch(prov, ref.User, ref.Repo, ref.Tag)
					if err == nil && cat != nil {
						for _, entry := range cat.Skills {
							suggestions = append(suggestions, refStr+entry.ID)
						}
					}
				}
			}
		} else {
			lock, err := manifest.LoadLock()
			if err == nil && lock != nil {
				for source := range lock.Skills {
					name := manifest.SkillName(source)
					if strings.HasPrefix(name, toComplete) {
						suggestions = append(suggestions, name)
					}
				}
			}
		}

		return suggestions, cobra.ShellCompDirectiveNoFileComp
	},
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

	// Resolve skill path (via catalog.json if available)
	var overridePath string
	cat, err := catalog.Fetch(prov, ref.User, ref.Repo, ref.Tag)
	if err == nil && cat != nil {
		if entry := cat.Find(ref.Skill); entry != nil && entry.Path != "" {
			overridePath = entry.Path
		}
	}

	data, err := installer.FetchFile(cloneURL, repoURL, ref.Skill, ref.Tag, overridePath, "SKILL.md")
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
