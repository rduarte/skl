package cmd

import (
	"fmt"

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
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("skl version %s\n", Version))
}
