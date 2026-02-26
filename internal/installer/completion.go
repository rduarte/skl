package installer

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ConfigureCompletion installs the shell completion scripts for Bash or Zsh
// in standard user directories based on the current shell detection.
func ConfigureCompletion(binPath string) error {
	shell := detectShell()

	switch shell {
	case "bash":
		bashDir := filepath.Join(os.Getenv("HOME"), ".local/share/bash-completion/completions")
		if err := os.MkdirAll(bashDir, 0755); err == nil {
			bashPath := filepath.Join(bashDir, "skl")
			_ = writeCompletion(binPath, "bash", bashPath)
		}
	case "zsh":
		zshDir := filepath.Join(os.Getenv("HOME"), ".local/share/zsh/site-functions")
		if err := os.MkdirAll(zshDir, 0755); err == nil {
			zshPath := filepath.Join(zshDir, "_skl")
			_ = writeCompletion(binPath, "zsh", zshPath)
		}
	default:
		// If shell detection fails or is non-standard, we don't spam the user.
		// Standard completion scripts are still available via 'skl completion'.
	}

	return nil
}

func detectShell() string {
	shellPath := os.Getenv("SHELL")
	if strings.Contains(shellPath, "zsh") {
		return "zsh"
	}
	if strings.Contains(shellPath, "bash") {
		return "bash"
	}
	return ""
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
