package utils

import (
	"context"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/pterm/pterm"
	"golang.org/x/oauth2"
)

type RepData struct {
	Name string
	Url  string
	Size uint64
}

// Get repository list from an organization that has a specific prefix
func GetRepositoryList(org, rep string) []RepData {
	ghToken := GetAndSaveToken(false)

	s, _ := pterm.DefaultSpinner.Start("Searching for '" + rep)
	repositories, response, err := getRepoOrg(ghToken, org, rep)
	if err != nil {
		s.Fail(response.Request.URL.String() + " -> " + response.Status + "\n")
		os.Exit(1)
	}
	s.Success()
	return repositories
}

// Get all repositories from an organization that has a specific prefix
func getRepoOrg(ghToken string, org string, repoPrefix string) ([]RepData, *github.Response, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	reposData := make([]RepData, 0)

	nextPage := 1
	lastPage := int(^uint(0) >> 1)

	for nextPage <= lastPage {
		page := github.ListOptions{Page: nextPage}
		repoOpts := github.RepositoryListByOrgOptions{
			Type:        "all",
			ListOptions: page,
		}
		// lists the repositories for an organization
		repos, response, err := client.Repositories.ListByOrg(ctx, org, &repoOpts)

		if err != nil {
			return nil, response, err
		}
		lastPage = response.LastPage

		for _, r := range repos {
			url := github.Stringify(r.HTMLURL)
			url = strings.ReplaceAll(url, "\"", "")
			if strings.HasPrefix(*r.Name, repoPrefix) {
				n := RepData{Name: *r.Name, Url: url, Size: uint64(*r.Size)}
				reposData = append(reposData, n)
			}
		}
		nextPage++
	}
	return reposData, nil, nil
}
