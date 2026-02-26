package installer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const skillsDir = ".agent/skills"

// Install clones the given repo using sparse-checkout and copies only the
// skill subdirectory into .agent/skills/<skill> relative to the current
// working directory. If force is true, an existing skill is removed first.
func Install(cloneURL, repoURL, skill, tag, overridePath string, force bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório atual: %w", err)
	}

	destDir := filepath.Join(cwd, skillsDir, skill)

	// Check if skill already exists locally
	if _, err := os.Stat(destDir); err == nil {
		if !force {
			return fmt.Errorf("skill %q já existe em %s (remova manualmente para reinstalar)", skill, destDir)
		}
		// Force mode: remove existing skill
		if err := os.RemoveAll(destDir); err != nil {
			return fmt.Errorf("erro ao remover skill existente: %w", err)
		}
	}

	// Create temp dir for the sparse clone
	tmpDir, err := os.MkdirTemp("", "skl-clone-*")
	if err != nil {
		return fmt.Errorf("erro ao criar diretório temporário: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	fmt.Printf("⬇  Baixando skill %q...\n", skill)

	// Step 1: Clone with sparse-checkout (blob filter for speed, no file checkout yet)
	cloneArgs := []string{
		"clone",
		"--filter=blob:none",
		"--sparse",
		"--depth=1",
		"--no-checkout",
	}
	if tag != "" {
		cloneArgs = append(cloneArgs, "--branch", tag)
	}
	cloneArgs = append(cloneArgs, cloneURL, tmpDir)

	if stderr, err := runGitCapture(cloneArgs...); err != nil {
		return classifyCloneError(stderr, repoURL, tag)
	}

	// Step 2: Resolve the skill directory inside the repo
	skillRepoPath := overridePath
	if skillRepoPath == "" {
		// Try .agent/skills/<skill> first
		primaryPath := filepath.Join(".agent/skills", skill)
		if err := verifyPathExists(tmpDir, primaryPath, repoURL); err == nil {
			skillRepoPath = primaryPath
		} else {
			// Try skills/<skill> as fallback
			fallbackPath := filepath.Join("skills", skill)
			if err := verifyPathExists(tmpDir, fallbackPath, repoURL); err == nil {
				skillRepoPath = fallbackPath
			} else {
				// If both fail, return the primary error for clarity
				return err
			}
		}
	} else if strings.HasSuffix(skillRepoPath, "/SKILL.md") || skillRepoPath == "SKILL.md" {
		skillRepoPath = filepath.Dir(skillRepoPath)
	}

	// Step 3: Sparse-checkout only the skill directory and checkout
	if err := runGitC(tmpDir, "sparse-checkout", "set", skillRepoPath); err != nil {
		return fmt.Errorf("erro no sparse-checkout: %w", err)
	}
	if err := runGitC(tmpDir, "checkout"); err != nil {
		return fmt.Errorf("erro no checkout: %w", err)
	}

	// Step 4: Copy skill directory to .agent/skills/<skill>
	skillSrc := filepath.Join(tmpDir, skillRepoPath)
	if err := os.MkdirAll(filepath.Dir(destDir), 0o755); err != nil {
		return fmt.Errorf("erro ao criar diretório de destino: %w", err)
	}

	if err := copyDir(skillSrc, destDir); err != nil {
		return fmt.Errorf("erro ao copiar skill: %w", err)
	}

	fmt.Printf("✅ Skill %q instalada em %s (via %s)\n", skill, destDir, skillRepoPath)
	return nil
}

// FetchFile fetches a single file from a skill directory in a remote repo.
// It clones sparsely, reads the file, and cleans up the temp dir.
func FetchFile(cloneURL, repoURL, skill, tag, overridePath, filename string) ([]byte, error) {
	var skillRepoPath string

	tmpDir, err := os.MkdirTemp("", "skl-fetch-*")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar diretório temporário: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Clone (no checkout)
	cloneArgs := []string{
		"clone",
		"--filter=blob:none",
		"--sparse",
		"--depth=1",
		"--no-checkout",
	}
	if tag != "" {
		cloneArgs = append(cloneArgs, "--branch", tag)
	}
	cloneArgs = append(cloneArgs, cloneURL, tmpDir)

	if stderr, err := runGitCapture(cloneArgs...); err != nil {
		return nil, classifyCloneError(stderr, repoURL, tag)
	}

	// Resolve the skill directory inside the repo
	skillRepoPath = overridePath
	if skillRepoPath == "" {
		// Try .agent/skills/<skill> first
		primaryPath := filepath.Join(".agent/skills", skill)
		if err := verifyPathExists(tmpDir, primaryPath, repoURL); err == nil {
			skillRepoPath = primaryPath
		} else {
			// Try skills/<skill> as fallback
			fallbackPath := filepath.Join("skills", skill)
			if err := verifyPathExists(tmpDir, fallbackPath, repoURL); err == nil {
				skillRepoPath = fallbackPath
			} else {
				// If both fail, return the primary error for clarity
				return nil, err
			}
		}
	} else if strings.HasSuffix(skillRepoPath, "/SKILL.md") || skillRepoPath == "SKILL.md" {
		skillRepoPath = filepath.Dir(skillRepoPath)
	}

	// Sparse-checkout and checkout
	if err := runGitC(tmpDir, "sparse-checkout", "set", skillRepoPath); err != nil {
		return nil, fmt.Errorf("erro no sparse-checkout: %w", err)
	}
	if err := runGitC(tmpDir, "checkout"); err != nil {
		return nil, fmt.Errorf("erro no checkout: %w", err)
	}

	// Read the requested file
	filePath := filepath.Join(tmpDir, skillRepoPath, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("arquivo %q não encontrado na skill %q", filename, skill)
		}
		return nil, fmt.Errorf("erro ao ler %s: %w", filename, err)
	}

	return data, nil
}

// verifyPathExists uses "git ls-tree" to check if a path exists in the repo
// tree before attempting sparse-checkout. This gives a clear error early.
func verifyPathExists(repoDir, path, repoURL string) error {
	cmd := exec.Command("git", "-C", repoDir, "ls-tree", "HEAD", path)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("erro ao verificar se skill existe no repositório: %w", err)
	}
	if len(strings.TrimSpace(string(out))) == 0 {
		return fmt.Errorf(
			"skill não encontrada no repositório\n"+
				"  Caminho esperado: %s\n\n"+
				"  Verifique se o diretório existe no repositório remoto: %s",
			path, repoURL,
		)
	}
	return nil
}

// runGitCapture executes a git command silently, capturing stderr for error analysis.
func runGitCapture(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var stderr bytes.Buffer
	cmd.Stdout = nil
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stderr.String(), err
}

// runGit executes a git command silently (output suppressed).
func runGit(args ...string) error {
	_, err := runGitCapture(args...)
	return err
}

// classifyCloneError inspects git stderr and returns a user-friendly error.
func classifyCloneError(stderr, repoURL, tag string) error {
	low := strings.ToLower(stderr)

	switch {
	case strings.Contains(low, "not found") ||
		strings.Contains(low, "does not exist") ||
		strings.Contains(low, "not exist") ||
		strings.Contains(low, "repository not found"):
		return fmt.Errorf(
			"repositório não encontrado\n\n"+
				"  Verifique se o repositório existe e se você tem acesso: %s",
			repoURL,
		)

	case tag != "" && (strings.Contains(low, "not a valid ref") ||
		strings.Contains(low, "remote branch") ||
		strings.Contains(low, "not found in upstream")):
		return fmt.Errorf(
			"tag %q não encontrada no repositório\n\n"+
				"  Verifique as tags disponíveis em: %s",
			tag, repoURL,
		)

	case strings.Contains(low, "permission denied") ||
		strings.Contains(low, "could not read from remote"):
		return fmt.Errorf(
			"acesso negado ao repositório\n\n"+
				"  Verifique suas credenciais SSH e acesso ao repositório: %s",
			repoURL,
		)

	default:
		return fmt.Errorf(
			"erro ao clonar repositório\n\n"+
				"  Repositório: %s\n"+
				"  Detalhe: %s",
			repoURL, strings.TrimSpace(stderr),
		)
	}
}

// runGitC executes a git command inside a specific directory.
func runGitC(dir string, args ...string) error {
	fullArgs := append([]string{"-C", dir}, args...)
	return runGit(fullArgs...)
}

// copyDir recursively copies src directory to dst.
func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		// Skip .git directory
		if entry.Name() == ".git" {
			continue
		}

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// ResolveRef uses "git ls-remote" to find the exact commit hash for a given ref (branch, tag or *).
func ResolveRef(cloneURL, gitRef string) (string, error) {
	args := []string{"ls-remote", cloneURL}

	target := gitRef
	if target == "" || target == "*" {
		target = "HEAD"
	}

	args = append(args, target)

	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("erro ao resolver referência remota: %v (detalhe: %s)", err, strings.TrimSpace(stderr.String()))
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return "", fmt.Errorf("referência %q não encontrada no repositório remoto", target)
	}

	// Format is <hash>\t<ref>
	parts := strings.Split(lines[0], "\t")
	if len(parts) < 1 {
		return "", fmt.Errorf("formato de resposta do git inválido ao resolver ref")
	}

	return parts[0], nil
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// List returns the names of currently installed skills.
func List() ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(cwd, skillsDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var skills []string
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			skills = append(skills, e.Name())
		}
	}
	return skills, nil
}
