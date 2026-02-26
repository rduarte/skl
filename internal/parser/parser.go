package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// SkillRef holds all the parsed components of a skill reference.
type SkillRef struct {
	Provider string // e.g. "github", "bitbucket"
	User     string // e.g. "empresa"
	Repo     string // e.g. "repo-skills"
	Skill    string // e.g. "data-analyzer"
	Tag      string // e.g. "v1.2.0" (empty if not specified)
}

// RepoRef holds parsed components of a repository reference.
type RepoRef struct {
	Provider string
	User     string
	Repo     string
	Tag      string
}

// pattern matches: <provider>@<user>/<repo>/<skill>[:tag]
var pattern = regexp.MustCompile(
	`^([a-zA-Z0-9-]+)@([a-zA-Z0-9._-]+)/([a-zA-Z0-9._-]+)/([a-zA-Z0-9._-]+)(?::([a-zA-Z0-9._-]+))?$`,
)

// repoPattern matches: <provider>@<user>/<repo>[:tag]
var repoPattern = regexp.MustCompile(
	`^([a-zA-Z0-9-]+)@([a-zA-Z0-9._-]+)/([a-zA-Z0-9._-]+)(?::([a-zA-Z0-9._-]+))?$`,
)

// Parse takes a raw skill reference string and returns a SkillRef.
func Parse(raw string) (*SkillRef, error) {
	raw = strings.TrimSuffix(raw, "/")

	// Special case: local@skill-name
	if strings.HasPrefix(raw, "local@") {
		parts := strings.Split(raw, "@")
		if len(parts) == 2 && parts[1] != "" {
			return &SkillRef{
				Provider: "local",
				Skill:    parts[1],
			}, nil
		}
	}

	matches := pattern.FindStringSubmatch(raw)
	if matches == nil {
		return nil, fmt.Errorf(
			"formato inválido: %q\nFormato esperado: <provider>@<user>/<repo>/<skill>[:tag] ou local@<skill>\nExemplo: github@empresa/repo-skills/data-analyzer:v1.2.0 ou local@my-skill",
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

// ParseRepo takes a raw repository reference string and returns a RepoRef.
func ParseRepo(raw string) (*RepoRef, error) {
	raw = strings.TrimSuffix(raw, "/")
	matches := repoPattern.FindStringSubmatch(raw)
	if matches == nil {
		return nil, fmt.Errorf(
			"formato de repositório inválido: %q\nFormato esperado: <provider>@<user>/<repo>[:tag]\nExemplo: github@empresa/repo-skills",
			raw,
		)
	}

	return &RepoRef{
		Provider: matches[1],
		User:     matches[2],
		Repo:     matches[3],
		Tag:      matches[4],
	}, nil
}

// String returns a human-readable representation of the SkillRef.
func (r *SkillRef) String() string {
	s := fmt.Sprintf("%s@%s/%s/%s", r.Provider, r.User, r.Repo, r.Skill)
	if r.Tag != "" {
		s += ":" + r.Tag
	}
	return s
}
