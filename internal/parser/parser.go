package parser

import (
	"fmt"
	"regexp"
)

// SkillRef holds all the parsed components of a skill reference.
type SkillRef struct {
	Provider string // e.g. "github", "bitbucket"
	User     string // e.g. "empresa"
	Repo     string // e.g. "repo-skills"
	Skill    string // e.g. "data-analyzer"
	Tag      string // e.g. "v1.2.0" (empty if not specified)
}

// pattern matches: <provider>@<user>/<repo>/<skill>[:tag]
// provider: alphanumeric + hyphens
// user:     alphanumeric + hyphens + dots
// repo:     alphanumeric + hyphens + dots
// skill:    alphanumeric + hyphens + dots + underscores
// tag:      alphanumeric + hyphens + dots (optional, prefixed by ":")
var pattern = regexp.MustCompile(
	`^([a-zA-Z0-9-]+)@([a-zA-Z0-9._-]+)/([a-zA-Z0-9._-]+)/([a-zA-Z0-9._-]+)(?::([a-zA-Z0-9._-]+))?$`,
)

// Parse takes a raw skill reference string and returns a SkillRef.
func Parse(raw string) (*SkillRef, error) {
	matches := pattern.FindStringSubmatch(raw)
	if matches == nil {
		return nil, fmt.Errorf(
			"formato inválido: %q\nFormato esperado: <provider>@<user>/<repo>/<skill>[:tag]\nExemplo: github@empresa/repo-skills/data-analyzer:v1.2.0",
			raw,
		)
	}

	ref := &SkillRef{
		Provider: matches[1],
		User:     matches[2],
		Repo:     matches[3],
		Skill:    matches[4],
		Tag:      matches[5], // empty string if not captured
	}

	return ref, nil
}

// String returns a human-readable representation of the SkillRef.
func (r *SkillRef) String() string {
	s := fmt.Sprintf("%s@%s/%s/%s", r.Provider, r.User, r.Repo, r.Skill)
	if r.Tag != "" {
		s += ":" + r.Tag
	}
	return s
}
