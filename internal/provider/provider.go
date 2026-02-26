package provider

import "fmt"

// Provider knows how to build a Git clone URL for a specific hosting service.
type Provider interface {
	// Name returns the provider identifier (e.g. "github").
	Name() string

	// CloneURL returns the SSH clone URL for the given user/repo.
	CloneURL(user, repo string) string

	// RepoURL returns the browsable HTTPS URL for the given user/repo.
	RepoURL(user, repo string) string
}

// registry holds all known providers.
var registry = map[string]Provider{
	"github":    GitHub{},
	"bitbucket": Bitbucket{},
}

// New returns a Provider for the given name, or an error if unsupported.
func New(name string) (Provider, error) {
	p, ok := registry[name]
	if !ok {
		supported := make([]string, 0, len(registry))
		for k := range registry {
			supported = append(supported, k)
		}
		return nil, fmt.Errorf("provider %q não suportado (disponíveis: %v)", name, supported)
	}
	return p, nil
}

// Register adds a new provider to the registry. Useful for testing or plugins.
func Register(name string, p Provider) {
	registry[name] = p
}
