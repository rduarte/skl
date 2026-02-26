package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rduarte/skl/internal/catalog"
	"github.com/rduarte/skl/internal/parser"
	"github.com/rduarte/skl/internal/provider"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list <provider>@<user>/<repo>",
	Short: "Lista todas as skills de um repositório que possui um catalog.json",
	Long: `Busca o arquivo catalog.json no repositório indicado e lista todas as skills
disponíveis de forma organizada.

Exemplo:
  skl list github@rmyndharis/antigravity-skills`,
	Args: cobra.ExactArgs(1),
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	refStr := args[0]
	ref, err := parser.ParseRepo(refStr)
	if err != nil {
		return err
	}

	prov, err := provider.New(ref.Provider)
	if err != nil {
		return err
	}

	fmt.Printf("🔍 Buscando catálogo em %s/%s...\n", ref.User, ref.Repo)
	cat, err := catalog.Fetch(prov, ref.User, ref.Repo, ref.Tag)
	if err != nil {
		repoURL := prov.RepoURL(ref.User, ref.Repo)
		fmt.Printf("\n⚠  Este repositório não possui um catálogo organizado (catalog.json).\n")
		fmt.Printf("Sugerimos explorar o conteúdo manualmente através da URL:\n%s\n", repoURL)
		return nil
	}

	if len(cat.Skills) == 0 {
		fmt.Println("ℹ️  O catálogo está vazio.")
		return nil
	}

	fmt.Printf("\n📚 Skills encontradas em %s/%s (%d total):\n\n", ref.User, ref.Repo, len(cat.Skills))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SKILL ID\tCATEGORIA\tDESCRIÇÃO")
	fmt.Fprintln(w, "--------\t---------\t---------")

	for _, s := range cat.Skills {
		desc := s.Description
		if len(desc) > 60 {
			desc = desc[:57] + "..."
		}
		category := s.Category
		if category == "" {
			category = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", s.ID, category, desc)
	}
	w.Flush()

	fmt.Printf("\nPara instalar uma skill, use:\n  skl install %s/<skill-id>\n", refStr)

	return nil
}
