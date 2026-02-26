package cmd

import (
	"fmt"

	"github.com/rduarte/skl/internal/installer"
	"github.com/rduarte/skl/internal/manifest"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Indexa folders locais em .agent/skills como skills gerenciadas",
	Long: `Verifica subdiretórios em .agent/skills/ que não estão no manifesto 
e os adiciona como skills locais (local@nome-da-skill).`,
	RunE: runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetup(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Buscando skills locais não indexadas em .agent/skills/...")

	// 1. Get folder names from .agent/skills
	folders, err := installer.List()
	if err != nil {
		return fmt.Errorf("erro ao listar diretório de skills: %w", err)
	}

	if len(folders) == 0 {
		fmt.Println("✅ Nenhuma pasta encontrada em .agent/skills/.")
		return nil
	}

	// 2. Load current manifesto
	mf, err := manifest.Load()
	if err != nil {
		return fmt.Errorf("erro ao carregar %s: %w", manifest.FileName, err)
	}

	// 3. Map tracked skills to their folder names
	tracked := make(map[string]bool)
	for source := range mf.Skills {
		tracked[manifest.SkillName(source)] = true
	}

	addedCount := 0
	for _, folder := range folders {
		if !tracked[folder] {
			// Add as local@folder
			source := "local@" + folder
			fmt.Printf("➕ Indexando skill local: %q\n", folder)
			mf.Skills[source] = "*"
			addedCount++
		}
	}

	if addedCount == 0 {
		fmt.Println("✅ Todas as skills locais já estão indexadas no manifesto.")
		return nil
	}

	// 4. Save manifest and lock
	if err := mf.Save(); err != nil {
		return err
	}
	if err := mf.SaveLock(); err != nil {
		return err
	}

	fmt.Printf("\n✨ %d nova(s) skill(s) adicionada(s) ao %s e %s\n", addedCount, manifest.FileName, manifest.LockFileName)

	// Ensure lock file is ignored in .gitignore
	_ = manifest.EnsureIgnoreLock()

	return nil
}
