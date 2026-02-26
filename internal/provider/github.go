package provider

import "fmt"

// GitHub implements Provider for github.com repositories.
type GitHub struct{}

func (GitHub) Name() string { return "github" }

func (GitHub) CloneURL(user, repo string) string {
	return fmt.Sprintf("git@github.com:%s/%s.git", user, repo)
}

func (GitHub) RepoURL(user, repo string) string {
	return fmt.Sprintf("https://github.com/%s/%s", user, repo)
}
