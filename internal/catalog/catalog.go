package catalog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rduarte/skl/internal/provider"
)

// Catalog represents the catalog.json structure.
type Catalog struct {
	Skills []SkillEntry `json:"skills"`
}

// SkillEntry represents a single skill entry in the catalog.
type SkillEntry struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Path        string   `json:"path"`
}

// Fetch fetch the catalog.json from the given repository using the provider's RawURL.
func Fetch(prov provider.Provider, user, repo, ref string) (*Catalog, error) {
	rawURL := prov.RawURL(user, repo, ref, "catalog.json")

	client := http.Client{
		Timeout: 2 * time.Second, // Short timeout for autocomplete
	}

	resp, err := client.Get(rawURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar catálogo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("catálogo não encontrado no repositório (status %d)", resp.StatusCode)
	}

	var cat Catalog
	if err := json.NewDecoder(resp.Body).Decode(&cat); err != nil {
		return nil, fmt.Errorf("erro ao ler catálogo: %w", err)
	}

	return &cat, nil
}

// Find finds a skill by its ID in the catalog.
func (c *Catalog) Find(id string) *SkillEntry {
	for _, s := range c.Skills {
		if s.ID == id {
			return &s
		}
	}
	return nil
}
