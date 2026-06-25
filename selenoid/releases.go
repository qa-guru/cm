package selenoid

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func getGithubReleaseAssetURL(owner, repo, version, osName, arch, githubBaseUrl string) (string, error) {
	ctx := context.Background()
	client := github.NewClient(nil)
	if githubBaseUrl != "" {
		u, err := url.Parse(githubBaseUrl)
		if err != nil {
			return "", fmt.Errorf("invalid Github base url [%s]: %v", githubBaseUrl, err)
		}
		client.BaseURL = u
	}

	var release *github.RepositoryRelease
	var err error
	if version != Latest {
		release, _, err = client.Repositories.GetReleaseByTag(ctx, owner, repo, version)
	} else {
		release, _, err = client.Repositories.GetLatestRelease(ctx, owner, repo)
	}
	if err != nil {
		return "", err
	}
	if release == nil {
		return "", fmt.Errorf("unknown release: %s", version)
	}

	title := cases.Title(language.AmericanEnglish)
	for _, asset := range release.Assets {
		assetName := *(asset.Name)
		if strings.Contains(assetName, osName) && strings.Contains(assetName, arch) {
			return *(asset.BrowserDownloadURL), nil
		}
	}
	return "", fmt.Errorf("%s binary for %s %s is not available for release %s", repo, title.String(osName), arch, version)
}
