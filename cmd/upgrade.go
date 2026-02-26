package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/rduarte/skl/internal/installer"
	"github.com/spf13/cobra"
)

const (
	repoOwner = "rduarte"
	repoName  = "skl"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Atualiza o skl para a última versão",
	Long:  `Verifica a última versão disponível no GitHub e atualiza o binário automaticamente.`,
	Args:  cobra.NoArgs,
	RunE:  runUpgrade,
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func runUpgrade(cmd *cobra.Command, args []string) error {
	fmt.Printf("📦 Versão atual: %s\n", Version)
	fmt.Println("🔍 Verificando última versão...")

	// 1. Fetch latest release from GitHub API
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("erro ao consultar GitHub API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API retornou status %d. Verifique: https://github.com/%s/%s/releases",
			resp.StatusCode, repoOwner, repoName)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("erro ao interpretar resposta da API: %w", err)
	}

	// 2. Compare versions
	if release.TagName == Version {
		fmt.Printf("✅ Você já está na versão mais recente (%s)\n", Version)
		return nil
	}

	fmt.Printf("⬆  Nova versão disponível: %s\n", release.TagName)

	// 3. Find the right asset for this OS/arch
	assetName := fmt.Sprintf("skl-%s-%s", runtime.GOOS, runtime.GOARCH)
	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("binário %q não encontrado na release %s.\n  Verifique: https://github.com/%s/%s/releases/tag/%s",
			assetName, release.TagName, repoOwner, repoName, release.TagName)
	}

	// 4. Download the new binary
	fmt.Printf("⬇  Baixando %s...\n", assetName)
	binResp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("erro ao baixar binário: %w", err)
	}
	defer binResp.Body.Close()

	if binResp.StatusCode != http.StatusOK {
		return fmt.Errorf("falha ao baixar binário (status %d)", binResp.StatusCode)
	}

	// 5. Write to temp file, then replace current executable
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("erro ao obter caminho do executável: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "skl-upgrade-*")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := io.Copy(tmpFile, binResp.Body); err != nil {
		tmpFile.Close()
		return fmt.Errorf("erro ao baixar binário: %w", err)
	}
	tmpFile.Close()

	// Make executable
	if err := os.Chmod(tmpPath, 0o755); err != nil {
		return fmt.Errorf("erro ao definir permissões: %w", err)
	}

	// Replace the current binary
	if err := os.Rename(tmpPath, execPath); err != nil {
		// Rename might fail across filesystems, try copy
		if err := copyBinaryFile(tmpPath, execPath); err != nil {
			return fmt.Errorf("erro ao substituir binário: %w\n  Tente manualmente: sudo mv %s %s", err, tmpPath, execPath)
		}
	}

	// 6. Update completions
	installer.ConfigureCompletion(execPath)

	fmt.Printf("✅ skl atualizado para %s\n", release.TagName)
	return nil
}

// copyBinaryFile copies src to dst (fallback when rename fails across filesystems).
func copyBinaryFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
