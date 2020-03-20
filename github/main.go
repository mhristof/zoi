package github

import (
	"context"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mhristof/zoi/ansible/galaxy"
	"github.com/mhristof/zoi/log"

	"github.com/google/go-github/github"
	"github.com/hashicorp/go-version"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type GitHub struct {
	ctx    context.Context
	client *github.Client
	repo   *github.RepositoriesService
}

func New() *GitHub {
	var g GitHub
	g.client, g.ctx = githubClient()
	g.repo = g.client.Repositories
	return &g
}

type HubConfig struct {
	Github []map[string]interface{} `yaml:"github.com"`
}

func hubToken() (string, error) {
	var config HubConfig
	configData, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".config/hub"))
	if err != nil {
		return "", err
	}
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return "", err
	}

	return config.Github[0]["oauth_token"].(string), nil
}

func githubToken() string {
	token, err := hubToken()
	if err == nil {
		return token
	}

	token, found := os.LookupEnv("GITHUB_TOKEN")

	if found {
		return token
	}

	panic("Error, could not find GH token in the usual places")
}

func githubClient() (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken()},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), ctx
}

func extractUserRepoFromSrc(src string) (string, string) {
	log.WithFields(log.Fields{
		"src": src,
	}).Debug("Extracting repo")

	url, err := url.Parse(src)
	if strings.HasPrefix(src, "https://") && err == nil {
		// this is your normal github url
		// https://github.com/user/repo
		urlParts := strings.Split(url.Path, "/")
		return urlParts[1], urlParts[2]
	}

	parts := strings.Split(src, ".")
	if len(parts) == 2 {
		// this is a 'user.role' source
		_, user, role := galaxy.FindRoleURL(parts[0], parts[1])
		return user, role
	}

	log.WithFields(log.Fields{
		"src": src,
	}).Panic("Cannot handle source")
	return "", ""
}

func (g GitHub) LatestTag(src string) string {
	log.WithFields(log.Fields{
		"src": src,
	}).Debug("Retrieving tags")
	tags, err := g.Tags(extractUserRepoFromSrc(src))
	if err != nil {
		log.WithFields(log.Fields{
			"src": src,
		}).Warning("Could not retrieve tags")
		return ""
	}

	sort.Sort(ByVersionDesc(tags))
	if len(tags) == 0 {
		return ""
	}
	return *tags[0].Name
}

func (g GitHub) LatestBranchCommit(repo string) string {
	log.WithFields(log.Fields{
		"repo": repo,
	}).Debug("Retrieving latest commit ids")

	commits := g.Commits(extractUserRepoFromSrc(repo))
	if commits == nil {
		log.WithFields(log.Fields{
			"repo": repo,
		}).Warning("Could not retrieve commits")
		return ""
	}

	sort.Sort(ByCommitTimeDesc(commits))

	return *commits[0].SHA
}

type ByCommitTimeDesc []*github.RepositoryCommit

func (a ByCommitTimeDesc) Len() int      { return len(a) }
func (a ByCommitTimeDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCommitTimeDesc) Less(i, j int) bool {
	return a[i].Commit.Author.Date.After(*a[j].Commit.Author.Date)
}

func (g GitHub) Commits(owner, repo string) []*github.RepositoryCommit {
	log.WithFields(log.Fields{
		"owner": owner,
		"repo":  repo,
	}).Debug("Retrieving all commits")

	var allCommits []*github.RepositoryCommit

	opts := github.CommitsListOptions{}

	for {
		commits, resp, err := g.client.Repositories.ListCommits(g.ctx, owner, repo, &opts)
		if err != nil {
			return nil
		}
		allCommits = append(allCommits, commits...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allCommits
}

type ByVersionDesc []*github.RepositoryTag

func (a ByVersionDesc) Len() int      { return len(a) }
func (a ByVersionDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByVersionDesc) Less(i, j int) bool {
	v1, err := version.NewVersion(*a[i].Name)
	if err != nil {
		return false
	}

	v2, err := version.NewVersion(*a[j].Name)
	if err != nil {
		return true
	}

	return v1.GreaterThan(v2)
}

func (g GitHub) Tags(owner string, repo string) ([]*github.RepositoryTag, error) {
	log.WithFields(log.Fields{
		"owner": owner,
		"repo":  repo,
	}).Debug("Retrieving tags")
	opt := &github.ListOptions{
		PerPage: 100,
	}

	var allTags []*github.RepositoryTag

	for {
		tags, resp, err := g.client.Repositories.ListTags(g.ctx, owner, repo, opt)
		if err != nil {
			return nil, err
		}
		allTags = append(allTags, tags...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allTags, nil
}
