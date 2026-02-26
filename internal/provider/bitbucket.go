package provider

import "fmt"

// Bitbucket implements Provider for bitbucket.org repositories.
type Bitbucket struct{}

func (Bitbucket) Name() string { return "bitbucket" }

func (Bitbucket) CloneURL(user, repo string) string {
	return fmt.Sprintf("git@bitbucket.org:%s/%s.git", user, repo)
}

func (Bitbucket) RepoURL(user, repo string) string {
	return fmt.Sprintf("https://bitbucket.org/%s/%s", user, repo)
}

func (Bitbucket) RawURL(user, repo, ref, path string) string {
	if ref == "" {
		ref = "main"
	}
	return fmt.Sprintf("https://bitbucket.org/%s/%s/raw/%s/%s", user, repo, ref, path)
}
