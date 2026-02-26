package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rduarte/skl/internal/manifest"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <skill-name>",
	Short: "Remove uma skill instalada",
	Long: `Remove uma skill do diretório .agent/skills/ e do sklfile.json.

Exemplo:
  skl remove 1doc-api-expert`,
	Args: cobra.ExactArgs(1),
	RunE: runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	skill := args[0]

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório atual: %w", err)
	}

	skillDir := filepath.Join(cwd, ".agent", "skills", skill)
	removed := false

	// 1. Remove skill directory if it exists
	if _, err := os.Stat(skillDir); err == nil {
		if err := os.RemoveAll(skillDir); err != nil {
			return fmt.Errorf("erro ao remover diretório: %w", err)
		}
		fmt.Printf("🗑️  Diretório removido: .agent/skills/%s\n", skill)
		removed = true
	}

	// 2. Remove from sklfile.json if listed
	mf, err := manifest.Load()
	if err != nil {
		return fmt.Errorf("erro ao carregar %s: %w", manifest.FileName, err)
	}

	// Find the key that ends with /<skill>
	var matchedKey string
	for source := range mf.Skills {
		parts := strings.Split(source, "/")
		if len(parts) > 0 && parts[len(parts)-1] == skill {
			matchedKey = source
			break
		}
	}

	if matchedKey != "" {
		delete(mf.Skills, matchedKey)
		if err := mf.Save(); err != nil {
			return fmt.Errorf("erro ao atualizar %s: %w", manifest.FileName, err)
		}
		fmt.Printf("📋 Removida do %s: %s\n", manifest.FileName, matchedKey)
		removed = true
	}

	if !removed {
		return fmt.Errorf("skill %q não encontrada (nem instalada, nem no %s)", skill, manifest.FileName)
	}

	fmt.Printf("✅ Skill %q removida\n", skill)
	return nil
}
