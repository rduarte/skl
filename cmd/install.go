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

		// Only autocomplete if we have at least provider@user/repo
		if !strings.Contains(toComplete, "@") {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

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

		if ref == nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
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

var forceInstall bool

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVarP(&forceInstall, "force", "f", false, "Sobrescreve a skill se ela já estiver instalada")
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

	// 4. Resolve skill path (via catalog.json if available)
	var overridePath string
	cat, err := catalog.Fetch(prov, ref.User, ref.Repo, ref.Tag)
	if err == nil && cat != nil {
		if entry := cat.Find(ref.Skill); entry != nil && entry.Path != "" {
			overridePath = entry.Path
			fmt.Printf("📖 Skill localizada via catálogo: %s\n", overridePath)
		}
	}

	// 5. Install the skill (force=false: don't overwrite existing)
	if err := installer.Install(cloneURL, repoURL, ref.Skill, ref.Tag, overridePath, forceInstall); err != nil {
		return err
	}

	// 6. Register in sklfile.json
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

	// 7. Update sklfile.lock with the exact commit hash
	lock, err := manifest.LoadLock()
	if err != nil {
		return fmt.Errorf("erro ao carregar %s: %w", manifest.LockFileName, err)
	}

	// Local skills don't have a remote hash
	if ref.Provider == "local" {
		lock.Skills[source] = "*"
	} else {
		hash, err := installer.ResolveRef(cloneURL, ref.Tag)
		if err != nil {
			fmt.Printf("⚠️  Aviso: Não foi possível resolver o hash remoto para o %s: %v\n", manifest.LockFileName, err)
			hash = gitRef // Fallback to symbolic ref
		}
		lock.Skills[source] = hash
	}

	if err := lock.SaveLock(); err != nil {
		return fmt.Errorf("erro ao atualizar %s: %w", manifest.LockFileName, err)
	}

	fmt.Printf("🔒 Skill bloqueada com hash no %s\n", manifest.LockFileName)

	// Ensure lock file is ignored in .gitignore
	_ = manifest.EnsureIgnoreLock()

	return nil
}
