//
// main.go
// Copyright (C) 2020 mhristof <mhristof@Mikes-MBP>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/hashicorp/go-version"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type Requirement struct {
	Src     string `yaml:"src"`
	Version string `yaml:"version"`
}

type GitHub struct {
	ctx    context.Context
	client *github.Client
	repo   *github.RepositoriesService
}

func (g *GitHub) New() {
	g.client, g.ctx = githubClient()
	g.repo = g.client.Repositories
}

func main() {
	if len(os.Args) < 2 {
		panic("Error, expected one argument")
	}

	requirementsPath := os.Args[1]
	if _, err := os.Stat(requirementsPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("Error, file %s does not exist", requirementsPath))
	}

	requirements := loadRequirementsFile(requirementsPath)

	// dont forget to import "encoding/json"
	requirementsJSON, err := json.MarshalIndent(requirements, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(requirementsJSON))

	var gh GitHub
	gh.New()

	var latestRequirements []Requirement
	for _, requirement := range requirements {
		latest := gh.latestTag(requirement)
		latestRequirements = append(latestRequirements, *latest)
	}

	dataOut, err := yaml.Marshal(&latestRequirements)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("latest.yml", dataOut, 0644)
	if err != nil {
		panic(err)
	}

}

func (g GitHub) latestTag(requirement Requirement) *Requirement {
	url, err := url.Parse(requirement.Src)
	if err != nil {
		panic(err)
	}

	urlParts := strings.Split(url.Path, "/")
	tags, err := g.Tags(urlParts[1], urlParts[2])
	if err != nil {
		panic(err)
	}

	sort.Sort(ByVersionDesc(tags))

	fmt.Println("source: ", requirement.Src)
	newReq := requirement
	newReq.Version = *tags[0].Name
	return &newReq
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
	fmt.Println("retrieving releases for", owner, repo)
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

type HubConfig struct {
	Github []map[string]interface{} `yaml:"github.com"`
}

func hubToken() (string, error) {
	var config HubConfig
	configData, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".config/hub"))
	fmt.Println(string(configData))
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

func loadRequirementsFile(path string) []Requirement {
	requirementsData, err := ioutil.ReadFile(path)
	var requirementList []Requirement
	err = yaml.Unmarshal(requirementsData, &requirementList)
	if err != nil {
		panic(err)
	}

	return requirementList
}
