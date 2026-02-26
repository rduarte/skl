package cmd

import (
	"fmt"
	"strings"

	"github.com/rduarte/skl/internal/catalog"
	"github.com/rduarte/skl/internal/installer"
	"github.com/rduarte/skl/internal/manifest"
	"github.com/rduarte/skl/internal/parser"
	"github.com/rduarte/skl/internal/provider"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install <provider>@<user>/<repo>/<skill>[:tag]",
	Short: "Baixa e instala uma skill no projeto atual",
	Long: `Baixa uma skill de um repositório Git e a instala em .agent/skills/<skill>.

Exemplos:
  skl install github@empresa/repo-skills/data-analyzer:v1.2.0
  skl install bitbucket@servicos-1doc/1doc-apis/1doc-api-expert`,
	Args: cobra.ExactArgs(1),
	RunE: runInstall,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		// Only autocomplete if we have at least provider@user/repo/
		if !strings.Contains(toComplete, "@") || !strings.Contains(toComplete, "/") {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		parts := strings.Split(toComplete, "/")
		// We need at least provider@user/repo/ (which means parts length >= 2 and the last part might be the skill prefix)
		if len(parts) < 2 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		// Extract ref parts from the first segments
		refStr := strings.Join(parts[:len(parts)-1], "/") + "/"
		ref, err := parser.Parse(refStr + "dummy") // add dummy skill to satisfy parser
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		prov, err := provider.New(ref.Provider)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		cat, err := catalog.Fetch(prov, ref.User, ref.Repo, ref.Tag)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var suggestions []string
		for _, entry := range cat.Skills {
			suggestions = append(suggestions, refStr+entry.ID)
		}

		return suggestions, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	// 1. Parse the skill reference
	ref, err := parser.Parse(args[0])
	if err != nil {
		return err
	}

	// 2. Resolve the provider
	prov, err := provider.New(ref.Provider)
	if err != nil {
		return err
	}

	// 3. Build the clone URL and the browsable repo URL
	cloneURL := prov.CloneURL(ref.User, ref.Repo)
	repoURL := prov.RepoURL(ref.User, ref.Repo)

	fmt.Printf("🔗 Clone URL: %s\n", cloneURL)

	// 4. Install the skill (force=false: don't overwrite existing)
	if err := installer.Install(cloneURL, repoURL, ref.Skill, ref.Tag, false); err != nil {
		return err
	}

	// 5. Register in sklfile.json
	// Key: provider@user/repo/skill  Value: tag (or empty)
	source := fmt.Sprintf("%s@%s/%s/%s", ref.Provider, ref.User, ref.Repo, ref.Skill)
	mf, err := manifest.Load()
	if err != nil {
		return fmt.Errorf("erro ao carregar %s: %w", manifest.FileName, err)
	}

	gitRef := ref.Tag
	if gitRef == "" {
		gitRef = "*"
	}

	if err := mf.Add(source, gitRef); err != nil {
		return fmt.Errorf("erro ao registrar skill no %s: %w", manifest.FileName, err)
	}

	// 6. Update sklfile.lock
	if err := mf.SaveLock(); err != nil {
		return fmt.Errorf("erro ao atualizar %s: %w", manifest.LockFileName, err)
	}

	fmt.Printf("📋 Skill registrada no %s e %s\n", manifest.FileName, manifest.LockFileName)
	return nil
}
