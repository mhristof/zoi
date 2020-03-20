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

type Requirement struct {
	Src     string `taml:"role,omitempty"`
	Role    string `yaml:"role,omitempty"`
	Version string `yaml:"version,omitempty"`
	Name    string `yaml:"name,omitempty"`
}

type Requirements []Requirement

type RoleListRequirement struct {
	Role map[string]Requirement
}

func (r *Requirement) updateSrc() {
	log.WithFields(log.Fields{
		"r": fmt.Sprintf("%+v", *r),
	}).Debug("Updating requirement fields")

	req := Requirement{}

	if r.Src != "" {
		r.Src = sanitiseGitURL(r.Src)
		return
	}

	if r.Role != "" {
		split := strings.Split(r.Role, ".")
		r.Src, _, _ = galaxy.FindRoleURL(split[0], split[1])
	} else if r.Name != "" {
		split := strings.Split(r.Name, ".")
		r.Src, _, _ = galaxy.FindRoleURL(split[0], split[1])
	}

	if req.Src == "" {
		return
	}

	r.Src = sanitiseGitURL(r.Src)
}

func sanitiseGitURL(url string) string {
	return strings.TrimPrefix(
		strings.TrimSuffix(url, ".git"),
		"git+")
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
	r.LoadBytes(requirementsData)
}

func (r *Requirements) LoadBytes(requirementsData []byte) {
	iface := loadRequirementsYAML(requirementsData)

	for _, item := range iface {
		itemJSON, err := json.MarshalIndent(item, "", "    ")
		if err != nil {
			panic(err)
		}

		var req Requirement
		err = json.Unmarshal(itemJSON, &req)
		if err != nil {
			continue
		}

		req.updateSrc()
		*r = append(*r, req)
	}
}

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

func (r *Requirements) Update() *Requirements {
	gh := github.New()

	var latestRequirements Requirements
	for _, requirement := range *r {
		log.WithFields(log.Fields{
			"requirement": requirement,
		}).Debug("Handling requirement")
		latest := gh.LatestTag(requirement.Src)
		if latest == "" {
			latest = gh.LatestBranchCommit(requirement.Src)
		}
		latestRequirements = append(latestRequirements, Requirement{
			Src:     requirement.Src,
			Version: latest,
		})
	}

	return &latestRequirements
}
