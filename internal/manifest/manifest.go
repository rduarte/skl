package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const FileName = "sklfile.json"

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

func filePath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("erro ao obter diretório atual: %w", err)
	}
	return filepath.Join(cwd, FileName), nil
}
