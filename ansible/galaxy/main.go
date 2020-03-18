package galaxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mhristof/zoi/log"
)

type GalaxyRolesResponse struct {
	Count        int
	Next         string
	NextLink     string
	Previous     string
	PreviousLink string
	Results      []AnsibleGalaxyRole
}

type SummaryFields struct {
	ContentType       map[string]interface{}
	Dependencies      []interface{}
	Namespace         map[string]interface{}
	Platforms         []interface{}
	ProviderNamespace map[string]interface{} `json:"provider_namespace"`
	Repository        map[string]interface{}
	Tags              []string
	Versions          []interface{}
	Videos            []interface{}
}

type AnsibleGalaxyRole struct {
	ID                int    `json:"id"`
	URL               string `json:"url"`
	Related           map[string]interface{}
	SummaryFields     SummaryFields
	Created           string
	Modified          string
	Name              string
	RoleType          string `json:"role_type"`
	IsValid           bool   `json:"is_valid"`
	MinAnsibleVersion string `json:"min_ansible_version"`
	License           string
	Company           string
	Description       string
	TravisStatusURL   string `json:"travis_status_url"`
	DownloadCount     int    `json:"download_count"`
	Imported          string
	Active            bool
	GithubUser        string `json:"github_user"`
	GithubRepo        string `json:"github_repo"`
	GithubServer      string `json:"github_server"`
	GithubBranch      string `json:"github_branch"`
	StargazersCount   int    `json:"stargazers_count"`
	ForksCount        int    `json:"forks_count"`
	OpenIssuesCount   int    `json:"open_issues_count"`
	Commit            string
	CommitMessage     string `json:"commit_message"`
	CommitURL         string `json:"commit_url"`
	IssueTrackerURL   string `json:"issue_tracker_url"`
}

func FindRoleURL(user, role string) (string, string, string) {
	url := fmt.Sprintf("https://galaxy.ansible.com/api/v1/roles/?owner__username=%s&name=%s", user, role)
	resp, err := http.Get(url)
	if err != nil {
		log.WithFields(log.Fields{
			"user": user,
			"role": role,
			"url":  url,
			"err":  err,
		}).Panic("Error while retrieving role from galaxy")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var gResp GalaxyRolesResponse
	err = json.Unmarshal(body, &gResp)
	if err != nil {
		log.WithFields(log.Fields{
			"user": user,
			"role": role,
			"url":  url,
			"err":  err,
		}).Panic("Error while unmarshaling galaxy response")
	}

	if gResp.Count != 1 {
		log.WithFields(log.Fields{
			"user":  user,
			"role":  role,
			"url":   url,
			"count": gResp.Count,
		}).Panic("Found more than 1 roles, im confused")
	}

	gRole := gResp.Results[0]

	ret := strings.Join([]string{
		gRole.GithubServer,
		gRole.GithubUser,
		gRole.GithubRepo,
	}, "/")
	return ret, gRole.GithubUser, gRole.GithubRepo
}
