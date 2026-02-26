package installer

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ConfigureCompletion installs the shell completion scripts for Bash and Zsh
// in standard user directories. These directories are usually loaded automatically
// by the shell if bash-completion or a standard Zsh setup is present.
func ConfigureCompletion(binPath string) error {
	// 1. Bash completion
	bashDir := filepath.Join(os.Getenv("HOME"), ".local/share/bash-completion/completions")
	if err := os.MkdirAll(bashDir, 0755); err == nil {
		bashPath := filepath.Join(bashDir, "skl")
		if err := writeCompletion(binPath, "bash", bashPath); err != nil {
			fmt.Printf("⚠  Não foi possível instalar autocompletar para Bash: %v\n", err)
		} else {
			fmt.Printf("✅ Autocompletar para Bash instalado em %s\n", bashPath)
		}
	}

	// 2. Zsh completion
	// Standard user path for Zsh completions
	zshDir := filepath.Join(os.Getenv("HOME"), ".local/share/zsh/site-functions")
	if err := os.MkdirAll(zshDir, 0755); err == nil {
		zshPath := filepath.Join(zshDir, "_skl")
		if err := writeCompletion(binPath, "zsh", zshPath); err != nil {
			fmt.Printf("⚠  Não foi possível instalar autocompletar para Zsh: %v\n", err)
		} else {
			fmt.Printf("✅ Autocompletar para Zsh instalado em %s\n", zshPath)
			fmt.Println("   Note: Certifique-se que ~/.local/share/zsh/site-functions está no seu $fpath")
		}
	}

	return nil
}

func writeCompletion(binPath, shell, dst string) error {
	cmd := exec.Command(binPath, "completion", shell)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return err
	}

	return os.WriteFile(dst, out.Bytes(), 0644)
}

// ShellHelper returns a one-liner to be added to ~/.bashrc or ~/.zshrc if needed,
// but our goal is to use standard paths to avoid manual config.
func ShellHelper() {
	// Not used for now as we prefer the standard paths approach
}
