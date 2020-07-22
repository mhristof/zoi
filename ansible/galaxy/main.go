package galaxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mhristof/zoi/log"
)

type galaxyRolesResponse struct {
	Count        int
	Next         string
	NextLink     string
	Previous     string
	PreviousLink string
	Results      []ansibleGalaxyRole
}

type summaryFields struct {
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

type ansibleGalaxyRole struct {
	ID                int    `json:"id"`
	URL               string `json:"url"`
	Related           map[string]interface{}
	SummaryFields     summaryFields
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

func get(url string) (resp *http.Response, err error) {
	for i := 0; i < 5; i++ {
		resp, err = http.Get(url)
		if resp.StatusCode != 520 {
			return resp, err
		}
		// I think galaxy.ansible.com is rate limiting requests and sometimes
		// when i run the huge tests i get a 520 back.
		// <html>
		// <head>
		// <title>520 Origin Error</title>
		// </head>
		// <body bgcolor=\"white\">
		// 	<center><h1>520 Origin Error</h1></center>
		// <hr>
		// <center>cloudflare-nginx</center>
		// </body>
		// </html>
		log.WithFields(log.Fields{
			"i":               i,
			"resp.StatusCode": resp.StatusCode,
			"url":             url,
		}).Debug("Got a 520, sleeping it off")
		time.Sleep(100 * time.Millisecond)
	}

	return resp, err
}

// FindRoleURL Search ansible galaxy for a give user/role combination
// For found roles, the github URL, the github user and the github repository
// name will be returned, otherwise its empty strings
func FindRoleURL(user, role string) (string, string, string) {
	url := fmt.Sprintf("https://galaxy.ansible.com/api/v1/roles/?owner__username=%s&name=%s", user, role)

	log.WithFields(log.Fields{
		"user": user,
		"role": role,
		"url":  url,
	}).Debug("Querying ansible galaxy")

	resp, err := get(url)
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

	var gResp galaxyRolesResponse
	err = json.Unmarshal(body, &gResp)
	if err != nil {
		log.WithFields(log.Fields{
			"user": user,
			"role": role,
			"url":  url,
			"body": string(body),
			"err":  err,
		}).Panic("Error while unmarshaling galaxy response")
	}

	if gResp.Count != 1 {
		log.WithFields(log.Fields{
			"user":  user,
			"role":  role,
			"url":   url,
			"count": gResp.Count,
		}).Warning("Incorrect amount of ansible galaxy roles found")
		return "", "", ""
	}

	gRole := gResp.Results[0]

	if gRole.GithubServer == "" {
		gRole.GithubServer = "https://github.com"
	}

	ret := strings.Join([]string{
		gRole.GithubServer,
		gRole.GithubUser,
		gRole.GithubRepo,
	}, "/")

	log.WithFields(log.Fields{
		"ret": ret,
	}).Debug("Found role")

	return ret, gRole.GithubUser, gRole.GithubRepo
}
