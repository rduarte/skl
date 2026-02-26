package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rduarte/skl/internal/catalog"
	"github.com/rduarte/skl/internal/installer"
	"github.com/rduarte/skl/internal/manifest"
	"github.com/rduarte/skl/internal/parser"
	"github.com/rduarte/skl/internal/provider"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Sincroniza as skills com base no sklfile.json",
	Long: `Compara o sklfile.json (estado desejado) com o sklfile.lock (estado atual)
e executa as ações necessárias:

  • Nova skill no sklfile.json → instala
  • Skill removida do sklfile.json → remove
  • Versão alterada → remove e reinstala

Ao final, atualiza o sklfile.lock para refletir o estado atual.`,
	Args: cobra.NoArgs,
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	// 0. Check if manifest exists
	if _, err := os.Stat(manifest.FileName); os.IsNotExist(err) {
		fmt.Printf("⚠  Arquivo %s não encontrado neste diretório.\n", manifest.FileName)
		fmt.Println("   O comando 'update' requer um manifesto para sincronizar as skills.")
		return nil // Abort gracefully or return error? User said "informando que o arquivo não existe e que o comando foi abortado".
	}

	// Load desired state (sklfile.json)
	desired, err := manifest.Load()
	if err != nil {
		return err
	}

	// Load current state (sklfile.lock)
	locked, err := manifest.LoadLock()
	if err != nil {
		return err
	}

	// Compute diff
	toInstall, toRemove, toUpgrade := diffManifests(desired, locked)

	total := len(toInstall) + len(toRemove) + len(toUpgrade)
	if total == 0 {
		fmt.Println("✅ Tudo sincronizado — nenhuma alteração necessária")
		return nil
	}

	fmt.Printf("📋 Alterações detectadas:\n")
	if len(toInstall) > 0 {
		fmt.Printf("   + %d skill(s) para instalar\n", len(toInstall))
	}
	if len(toRemove) > 0 {
		fmt.Printf("   - %d skill(s) para remover\n", len(toRemove))
	}
	if len(toUpgrade) > 0 {
		fmt.Printf("   ↑ %d skill(s) para atualizar\n", len(toUpgrade))
	}
	fmt.Println()

	var errors []string
	success := 0

	// 1. Remove skills that were removed from sklfile.json
	for _, source := range toRemove {
		skill := manifest.SkillName(source)
		fmt.Printf("🗑️  Removendo %q...\n", skill)
		if err := removeSkillDir(skill); err != nil {
			errors = append(errors, fmt.Sprintf("  ✗ %s: %v", skill, err))
			continue
		}
		success++
	}

	// 2. Upgrade skills (remove old + install new)
	for _, source := range toUpgrade {
		skill := manifest.SkillName(source)
		oldRef := locked.Skills[source]
		newRef := desired.Skills[source]
		fmt.Printf("↑  Atualizando %q (%s → %s)...\n", skill, oldRef, newRef)

		// Remove old version
		if err := removeSkillDir(skill); err != nil {
			errors = append(errors, fmt.Sprintf("  ✗ %s: %v", skill, err))
			continue
		}

		// Install new version
		if err := installSkill(source, newRef); err != nil {
			errors = append(errors, fmt.Sprintf("  ✗ %s: %v", skill, err))
			continue
		}
		success++
		fmt.Println()
	}

	// 3. Install new skills
	for _, source := range toInstall {
		skill := manifest.SkillName(source)
		gitRef := desired.Skills[source]
		fmt.Printf("📦 Instalando %q...\n", skill)

		if err := installSkill(source, gitRef); err != nil {
			errors = append(errors, fmt.Sprintf("  ✗ %s: %v", skill, err))
			continue
		}
		success++
		fmt.Println()
	}

	// 4. Save lock file with current desired state
	if err := desired.SaveLock(); err != nil {
		return fmt.Errorf("erro ao salvar %s: %w", manifest.LockFileName, err)
	}

	// Summary
	fmt.Printf("📊 Resultado: %d/%d operação(ões) concluída(s)\n", success, total)
	fmt.Printf("🔒 %s atualizado\n", manifest.LockFileName)

	if len(errors) > 0 {
		fmt.Println("\n⚠  Erros:")
		for _, e := range errors {
			fmt.Println(e)
		}
	}

	return nil
}

// diffManifests compares desired (sklfile.json) vs locked (sklfile.lock)
// and returns lists of sources to install, remove, and upgrade.
func diffManifests(desired, locked *manifest.Manifest) (toInstall, toRemove, toUpgrade []string) {
	// New in desired, not in locked → install
	for source := range desired.Skills {
		if _, exists := locked.Skills[source]; !exists {
			toInstall = append(toInstall, source)
		}
	}

	// In locked, not in desired → remove
	for source := range locked.Skills {
		if _, exists := desired.Skills[source]; !exists {
			toRemove = append(toRemove, source)
		}
	}

	// In both, but different ref → upgrade
	for source, desiredRef := range desired.Skills {
		if lockedRef, exists := locked.Skills[source]; exists {
			if desiredRef != lockedRef {
				toUpgrade = append(toUpgrade, source)
			}
		}
	}

	return
}

// installSkill resolves provider and installs a skill.
func installSkill(source, gitRef string) error {
	fullRef := source
	if gitRef != "" && gitRef != "*" {
		fullRef += ":" + gitRef
	}

	ref, err := parser.Parse(fullRef)
	if err != nil {
		return err
	}

	prov, err := provider.New(ref.Provider)
	if err != nil {
		return err
	}

	cloneURL := prov.CloneURL(ref.User, ref.Repo)
	repoURL := prov.RepoURL(ref.User, ref.Repo)

	// Resolve skill path (via catalog.json if available)
	var overridePath string
	cat, err := catalog.Fetch(prov, ref.User, ref.Repo, ref.Tag)
	if err == nil && cat != nil {
		if entry := cat.Find(ref.Skill); entry != nil && entry.Path != "" {
			overridePath = entry.Path
		}
	}

	fmt.Printf("🔗 Clone URL: %s\n", cloneURL)
	return installer.Install(cloneURL, repoURL, ref.Skill, ref.Tag, overridePath, true)
}

// removeSkillDir removes the skill directory from .agent/skills/.
func removeSkillDir(skill string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	dir := filepath.Join(cwd, ".agent", "skills", skill)
	if _, err := os.Stat(dir); err == nil {
		return os.RemoveAll(dir)
	}
	return nil
}
