package cmd

import (
	"fmt"

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
