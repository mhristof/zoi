package gh

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

type Url struct {
	Host    string
	Owner   string
	Repo    string
	Release string
	Url     string
}

var (
	ErrorUrlTooShort     = errors.New("URL too short")
	ErrorWrongHost       = errors.New("URL host is wrong")
	ErrorNoReleases      = errors.New("No releases available")
	ErrorCannotHandleUrl = errors.New("Cannot handle the url")
)

func ParseGitUrl(url string) (*Url, error) {
	if !strings.HasPrefix(url, "git@github.com") {
		return nil, ErrorWrongHost
	}

	var ret = Url{
		Url: url,
	}
	release := strings.Split(url, "=")

	if len(release) == 2 {
		ret.Release = release[1]
	}

	host := strings.Split(url, ":")

	if len(host) != 2 {
		return nil, ErrorWrongHost
	}

	ret.Host = strings.Replace(host[0], "git@", "", -1)

	owner := strings.Split(host[1], "/")

	ret.Owner = owner[0]
	// remove everything after .git
	ret.Repo = owner[1][0:strings.Index(owner[1], ".git")]

	return &ret, nil
}

func ParseUrl(in string) (*Url, error) {
	url, err := ParseHttpUrl(in)
	if err == nil {
		return url, nil
	}

	url, err = ParseGitUrl(in)
	if err == nil {
		return url, nil
	}

	return nil, ErrorCannotHandleUrl
}

func ParseHttpUrl(url string) (*Url, error) {
	if !strings.HasPrefix(url, "https://github.com") {
		return nil, ErrorWrongHost
	}

	parts := strings.Split(url, "/")

	if len(parts) < 5 {
		return nil, ErrorUrlTooShort
	}

	return &Url{
		Host:    fmt.Sprintf("%s//%s", parts[0], parts[2]),
		Owner:   parts[3],
		Repo:    sanitiseRepo(parts[4]),
		Release: getRelease(parts...),
		Url:     url,
	}, nil
}

func sanitiseRepo(repo string) string {
	refPos := strings.Index(repo, "?ref")
	if refPos > 0 {
		repo = repo[0:refPos]
	}

	repo = strings.TrimSuffix(repo, ".git")
	return repo
}

func getRelease(parts ...string) string {

	if len(parts) > 7 && parts[5] == "releases" && parts[6] == "download" {
		return parts[7]
	}

	if len(parts) == 5 && strings.Contains(parts[4], "ref=") {
		return strings.Split(parts[4], "=")[1]
	}

	return ""
}

func (u *Url) NextRelease() (string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_READONLY_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.ListOptions{}

	releases, _, err := client.Repositories.ListReleases(ctx, u.Owner, u.Repo, opt)
	if err != nil {
		panic(err)
	}

	if len(releases) == 0 {
		return "", ErrorNoReleases
	}

	latest := releases[0]
	for _, release := range releases {
		this := semver.New(sanitiseRelease(*release.TagName))

		if semver.New(sanitiseRelease(*latest.TagName)).LessThan(*this) {
			latest = release
		}
	}

	return strings.Replace(u.Url, u.Release, *latest.TagName, -1), nil
}

func sanitiseRelease(tag string) string {
	return strings.TrimLeft(tag, "v")
}
