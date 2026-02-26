package cmd

import (
	"fmt"

	"github.com/rduarte/skl/internal/updater"
	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "skl",
	Short: "Gerenciador de skills para projetos de IA",
	Long: `skl é um gerenciador de dependências focado em instalar
"skills" (capacidades/ferramentas de IA) dentro de projetos locais.

Ele faz o download de skills armazenadas em repositórios Git
(GitHub, Bitbucket) e as organiza no diretório .agent/skills/.`,
	Version: Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip version check for specific commands
		skipped := []string{"upgrade", "completion", "help", "setup"}
		for _, s := range skipped {
			if cmd.Name() == s || (cmd.Parent() != nil && cmd.Parent().Name() == s) {
				return
			}
		}

		// Skip if Version is "dev"
		if Version == "dev" {
			return
		}

		// Perform check (ignore errors to not block user)
		latest, err := updater.CheckLatestVersion()
		if err == nil && latest != "" && latest != Version {
			fmt.Printf("\033[1;33m⚠  Nova versão do skl disponível: %s (atual: %s)\033[0m\n", latest, Version)
			fmt.Printf("\033[1;33m   Execute 'skl upgrade' para atualizar.\033[0m\n\n")
		}
	},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("skl version %s\n", Version))
}
