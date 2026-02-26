package cmd

import (
	"fmt"

	"github.com/rduarte/skl/internal/installer"
	"github.com/rduarte/skl/internal/manifest"
	"github.com/rduarte/skl/internal/parser"
	"github.com/rduarte/skl/internal/provider"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Instala/atualiza todas as skills listadas no sklfile.json",
	Long: `Lê o arquivo sklfile.json na raiz do projeto e instala ou atualiza
todas as skills listadas. Skills já existentes são removidas e
baixadas novamente para garantir a versão mais recente.`,
	Args: cobra.NoArgs,
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	mf, err := manifest.Load()
	if err != nil {
		return err
	}

	sources := mf.SortedSources()
	if len(sources) == 0 {
		fmt.Println("📋 Nenhuma skill encontrada no sklfile.json")
		return nil
	}

	fmt.Printf("📋 %d skill(s) encontrada(s) no sklfile.json\n\n", len(sources))

	var errors []string
	installed := 0

	for _, source := range sources {
		gitRef := mf.Skills[source]

		// Build the full reference for parsing: source[:tag]
		fullRef := source
		if gitRef != "" && gitRef != "*" {
			fullRef += ":" + gitRef
		}

		ref, err := parser.Parse(fullRef)
		if err != nil {
			errors = append(errors, fmt.Sprintf("  ✗ %s: %v", source, err))
			continue
		}

		prov, err := provider.New(ref.Provider)
		if err != nil {
			errors = append(errors, fmt.Sprintf("  ✗ %s: %v", source, err))
			continue
		}

		cloneURL := prov.CloneURL(ref.User, ref.Repo)
		repoURL := prov.RepoURL(ref.User, ref.Repo)

		fmt.Printf("🔗 Clone URL: %s\n", cloneURL)

		if err := installer.Install(cloneURL, repoURL, ref.Skill, ref.Tag, true); err != nil {
			errors = append(errors, fmt.Sprintf("  ✗ %s: %v", source, err))
			continue
		}

		installed++
		fmt.Println()
	}

	// Summary
	fmt.Printf("📊 Resultado: %d/%d skill(s) instalada(s)\n", installed, len(sources))

	if len(errors) > 0 {
		fmt.Println("\n⚠  Erros:")
		for _, e := range errors {
			fmt.Println(e)
		}
	}

	return nil
}
