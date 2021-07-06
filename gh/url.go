package gh

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/v33/github"
	"github.com/mhristof/zoi/log"
	"golang.org/x/oauth2"
)

type Url struct {
	Host    string
	Owner   string
	Repo    string
	Release string
	Url     string
	Token   string
}

var (
	ErrorURLTooShort      = errors.New("URL too short")
	ErrorWrongHost        = errors.New("URL host is wrong")
	ErrorNoReleases       = errors.New("no releases available")
	ErrorCannotHandleURL  = errors.New("cannot handle the url")
	ErrorNoTags           = errors.New("no tags available")
	ErrorReleaseNotInTags = errors.New("release string not a tag")
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

	ret.Host = strings.ReplaceAll(host[0], "git@", "")

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

	return nil, ErrorCannotHandleURL
}

func ParseHttpUrl(url string) (*Url, error) {
	if !strings.HasPrefix(url, "https://github.com") {
		return nil, ErrorWrongHost
	}

	parts := strings.Split(url, "/")

	if len(parts) < 5 {
		return nil, ErrorURLTooShort
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
	if u.Token == "" {
		log.WithFields(log.Fields{
			"url": u,
		}).Panic("url.Token not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: u.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	release, releaseErr := latestRelease(client, u.Owner, u.Repo)
	tag, tagErr := latestTag(client, u.Owner, u.Repo)

	if releaseErr == nil && tagErr == nil && tag != release {
		log.WithFields(log.Fields{
			"release": release,
			"tag":     tag,
			"u.Url":   u.Url,
		}).Warning("warning, latest tag doesnt match latest release")
	}

	if releaseErr != nil && tagErr == nil {
		release = tag
	}

	if releaseErr != nil && tagErr != nil {
		return "", ErrorCannotHandleURL
	}

	log.WithFields(log.Fields{
		"u.Url":     u.Url,
		"u.Release": u.Release,
		"release":   release,
	}).Debug("New release")

	return u.sanitize(release), nil
}

func latestTag(client *github.Client, owner, repo string) (string, error) {
	ctx := context.Background()
	opt := &github.ListOptions{}

	tags, _, err := client.Repositories.ListTags(ctx, owner, repo, opt)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Panic("Could not retrieve information from server")
	}

	if len(tags) == 0 {
		return "", ErrorNoTags
	}

	latest := tags[0]

	for _, tag := range tags {
		// skip prepreleases
		this, err := semver.NewVersion(sanitiseRelease(*tag.Name))
		if err != nil {
			continue
		}

		if this.PreRelease == "" {
			latest = tag

			break
		}
	}

	log.WithFields(log.Fields{
		"*latest.Name": *latest.Name,
	}).Debug("Latest release name")

	return *latest.Name, nil
}

func latestRelease(client *github.Client, owner, repo string) (string, error) {
	ctx := context.Background()
	opt := &github.ListOptions{}

	releases, _, err := client.Repositories.ListReleases(ctx, owner, repo, opt)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Panic("Could not retrieve information from server")
	}

	if len(releases) == 0 {
		return "", ErrorNoReleases
	}

	latest := releases[0]

	for _, release := range releases {
		this, err := semver.NewVersion(sanitiseRelease(*release.TagName))
		if err != nil {
			continue
		}

		if semver.New(sanitiseRelease(*latest.TagName)).LessThan(*this) {
			latest = release
		}
	}

	return *latest.TagName, nil
}

func sanitiseRelease(tag string) string {
	return strings.TrimLeft(tag, "v")
}

func (u *Url) sanitize(release string) string {
	releaseNew := strings.ReplaceAll(
		// replace all versions in string
		u.Url, u.Release, release,
	)

	currentReleaseWithoutV := strings.TrimPrefix(u.Release, "v")

	if !strings.Contains(release, currentReleaseWithoutV) || len(currentReleaseWithoutV) != 1 {
		// if the release is something like `v2` and the new release is
		// something like `v2.3.1`, then this replace should not happen as
		// it will result in a string like `v2.3.1.3.1`
		log.WithFields(log.Fields{
			"releaseNew":             releaseNew,
			"u":                      u,
			"release":                release,
			"currentReleaseWithoutV": currentReleaseWithoutV,
		}).Debug("Replacing all occurrences")

		releaseNew = strings.ReplaceAll(
			releaseNew,
			// replace version that might exist without the v prefix
			strings.TrimPrefix(u.Release, "v"),
			strings.TrimPrefix(release, "v"),
		)
	}

	return releaseNew
}
