package github

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

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

	parts := strings.Split(src, ".")
	if len(parts) == 2 {
		// this is a 'user.role' source
		return parts[0], parts[1]
	}

	url, err := url.Parse(src)
	if err == nil {
		// this is your normal github url
		// https://github.com/user/repo
		urlParts := strings.Split(url.Path, "/")
		return urlParts[0], urlParts[1]
	}

	log.WithFields(log.Fields{
		"src": src,
	}).Panic("Cannot handle source")
	return "", ""
}

func (g GitHub) LatestTag(src string) string {
	tags, err := g.Tags(extractUserRepoFromSrc(src))
	if err != nil {
		panic(err)
	}

	sort.Sort(ByVersionDesc(tags))

	fmt.Println("source: ", src)
	return *tags[0].Name
}

type ByVersionDesc []*github.RepositoryTag

func (a ByVersionDesc) Len() int      { return len(a) }
func (a ByVersionDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByVersionDesc) Less(i, j int) bool {
	v1, err := version.NewVersion(*a[i].Name)
	if err != nil {
		panic(err)
	}

	v2, err := version.NewVersion(*a[j].Name)
	if err != nil {
		panic(err)
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
