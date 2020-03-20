package ansible

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mhristof/zoi/ansible/galaxy"
	"github.com/mhristof/zoi/github"
	"github.com/mhristof/zoi/log"
	"gopkg.in/yaml.v3"
)

type requirement struct {
	Src     string `taml:"role,omitempty"`
	Role    string `yaml:"role,omitempty"`
	Version string `yaml:"version,omitempty"`
	Name    string `yaml:"name,omitempty"`
}

// Requirements Holds a list of requirements for a requirements.yml file
type Requirements []requirement

func (r *requirement) updateSrc() {
	log.WithFields(log.Fields{
		"r": fmt.Sprintf("%+v", *r),
	}).Debug("Updating requirement fields")

	if r.Src == "" {
		if r.Role != "" {
			r.Src = r.Role
		} else if r.Name != "" {
			r.Src = r.Name
		}
	}

	r.Src = sanitiseGitURL(r.Src)

	if !strings.HasPrefix(r.Src, "http") {
		split := strings.Split(r.Src, ".")
		r.Src, _, _ = galaxy.FindRoleURL(split[0], split[1])
	}
}

func sanitiseGitURL(url string) string {
	url = strings.TrimPrefix(url, "git+")

	if strings.HasPrefix(url, "http") {
		url = strings.TrimSuffix(url, ".git")
	}
	return url
}

func githubPreffix(user, role string) string {
	return ""
}

func loadRequirementsYAML(data []byte) []interface{} {
	var iface []interface{}
	err := yaml.Unmarshal(data, &iface)
	if err == nil {
		// yaml file is a lsit of requirements
		return iface
	}

	var rolesIface map[string]interface{}
	err = yaml.Unmarshal(data, &rolesIface)
	if err == nil {
		// yaml file is probably a dictionary with the requirements being under
		// "role" key
		return rolesIface["roles"].([]interface{})
	}

	log.Panic("Error, i dont know how to handle this requirements file")
	return nil
}

// LoadFromFile Loads a requirements.yml file into go
func (r *Requirements) LoadFromFile(path string) {
	log.WithFields(log.Fields{
		"path": path,
	}).Debug("Loading requirements file")

	requirementsData, err := ioutil.ReadFile(path)
	if err != nil {
		log.WithFields(log.Fields{
			"path": path,
		}).Panic("Error while reading file")
	}
	r.loadBytes(requirementsData)
}

func (r *Requirements) loadBytes(requirementsData []byte) {
	iface := loadRequirementsYAML(requirementsData)

	for _, item := range iface {
		itemJSON, err := json.MarshalIndent(item, "", "    ")
		if err != nil {
			panic(err)
		}

		var req requirement
		err = json.Unmarshal(itemJSON, &req)
		if err != nil {
			continue
		}

		req.updateSrc()
		if req.Src == "" {
			continue
		}
		*r = append(*r, req)
	}
}

// SaveToFile Dumps a golang representation of requirements to a file
func (r *Requirements) SaveToFile(path string) {
	dataOut, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path, dataOut, 0644)
	if err != nil {
		panic(err)
	}
}

// Update Updates the requirements to the latest and geatest
// In order of precedence, this fucntion will return
// 	1. the latest tag
//  2. the latest commit from master
func (r *Requirements) Update() *Requirements {
	gh := github.New()

	var latestRequirements Requirements
	for _, req := range *r {
		log.WithFields(log.Fields{
			"req": req,
		}).Debug("Handling requirement")
		latest := gh.LatestTag(req.Src)
		if latest == "" {
			latest = gh.LatestBranchCommit(req.Src)
		}
		latestRequirements = append(latestRequirements, requirement{
			Src:     req.Src,
			Version: latest,
		})
	}

	return &latestRequirements
}
