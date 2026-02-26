package cmd

import (
	"fmt"
	"os"

	"github.com/rduarte/skl/internal/installer"
	"github.com/spf13/cobra"
)

var rebuildCmd = &cobra.Command{
	Use:    "rebuild",
	Short:  "Reconstrói as configurações de autocompletar",
	Hidden: true, // Hidden because it's mostly for automation
	RunE: func(cmd *cobra.Command, args []string) error {
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("erro ao obter caminho do executável: %w", err)
		}
		return installer.ConfigureCompletion(execPath)
	},
}

func init() {
	rootCmd.AddCommand(rebuildCmd)
}
