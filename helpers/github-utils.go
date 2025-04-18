package helpers

import (
	"fmt"
	"strings"
)

type GithubUrl struct {
	Url string
	Parent string
	Repo   string
}

func ParseGithubUrl(githubUrl string) (GithubUrl, error) {
	// Check if the URL contains "github.com/"
	const githubPrefix = "https://github.com/"
	if !strings.HasPrefix(githubUrl, githubPrefix) {
		return GithubUrl{}, fmt.Errorf("invalid GitHub URL: %s", githubUrl)
	}

	// Remove the prefix and split the remaining part by "/"
	path := strings.TrimPrefix(githubUrl, githubPrefix)
	pathParts := strings.Split(strings.TrimSuffix(path, ".git"), "/")

	// Ensure the URL is in the correct format
	if len(pathParts) != 2 {
		return GithubUrl{}, fmt.Errorf("invalid GitHub URL format: %s", githubUrl)
	}

	return GithubUrl{
		Url: githubUrl,
		Parent: pathParts[0],
		Repo:   pathParts[1],
	}, nil
}