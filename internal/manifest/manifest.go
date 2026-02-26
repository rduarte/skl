package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const FileName = "sklfile.json"
const LockFileName = "sklfile.lock"

// Manifest represents the sklfile.json file.
// Keys are full skill references (e.g. "bitbucket@user/repo/skill"),
// values are the git ref (branch or tag, e.g. "master", "v1.2.0").
type Manifest struct {
	Skills map[string]string `json:"skills"`
}

// Load reads the manifest from sklfile.json in the current directory.
// Returns an empty manifest if the file does not exist.
func Load() (*Manifest, error) {
	path, err := filePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Manifest{Skills: make(map[string]string)}, nil
		}
		return nil, fmt.Errorf("erro ao ler %s: %w", FileName, err)
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("erro ao interpretar %s: %w", FileName, err)
	}

	if m.Skills == nil {
		m.Skills = make(map[string]string)
	}

	return &m, nil
}

// Save writes the manifest to sklfile.json in the current directory.
func (m *Manifest) Save() error {
	path, err := filePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar %s: %w", FileName, err)
	}

	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("erro ao gravar %s: %w", FileName, err)
	}

	return nil
}

// Add registers a skill in the manifest and saves it.
// source is the full reference (e.g. "bitbucket@user/repo/skill").
// ref is the git ref (branch or tag).
func (m *Manifest) Add(source, ref string) error {
	m.Skills[source] = ref
	return m.Save()
}

// SortedSources returns skill source keys sorted alphabetically.
func (m *Manifest) SortedSources() []string {
	sources := make([]string, 0, len(m.Skills))
	for s := range m.Skills {
		sources = append(sources, s)
	}
	sort.Strings(sources)
	return sources
}

// LoadLock reads the lock file (sklfile.lock). Returns empty manifest if absent.
func LoadLock() (*Manifest, error) {
	return loadFile(LockFileName)
}

// SaveLock writes the lock file (sklfile.lock).
func (m *Manifest) SaveLock() error {
	return m.saveFile(LockFileName)
}

// SkillName extracts the skill name (last path segment) from a source key.
// e.g. "bitbucket@servicos-1doc/1doc-apis/1doc-api-expert" → "1doc-api-expert"
func SkillName(source string) string {
	parts := strings.Split(source, "/")
	if len(parts) == 0 {
		return source
	}
	return parts[len(parts)-1]
}

func filePath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("erro ao obter diretório atual: %w", err)
	}
	return filepath.Join(cwd, FileName), nil
}

func loadFile(name string) (*Manifest, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter diretório atual: %w", err)
	}

	path := filepath.Join(cwd, name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Manifest{Skills: make(map[string]string)}, nil
		}
		return nil, fmt.Errorf("erro ao ler %s: %w", name, err)
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("erro ao interpretar %s: %w", name, err)
	}

	if m.Skills == nil {
		m.Skills = make(map[string]string)
	}

	return &m, nil
}

func (m *Manifest) saveFile(name string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório atual: %w", err)
	}

	path := filepath.Join(cwd, name)
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar %s: %w", name, err)
	}

	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("erro ao gravar %s: %w", name, err)
	}

	return nil
}

// EnsureIgnoreLock checks if .gitignore exists and contains LockFileName.
// If not, it adds it automatically.
func EnsureIgnoreLock() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	gitIgnorePath := filepath.Join(cwd, ".gitignore")
	if _, err := os.Stat(gitIgnorePath); os.IsNotExist(err) {
		return nil // No .gitignore, nothing to do
	}

	data, err := os.ReadFile(gitIgnorePath)
	if err != nil {
		return fmt.Errorf("erro ao ler .gitignore: %w", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	found := false
	for _, line := range lines {
		if strings.TrimSpace(line) == LockFileName {
			found = true
			break
		}
	}

	if !found {
		// Ensure it ends with newline if we're appending
		if len(content) > 0 && !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		content += LockFileName + "\n"
		if err := os.WriteFile(gitIgnorePath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("erro ao atualizar .gitignore: %w", err)
		}
	}

	return nil
}
