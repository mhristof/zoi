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
	Src     string `yaml:"src,omitempty"`
	Version string `yaml:"version,omitempty"`
}

type Requirements []Requirement

type RoleRequirement struct {
	Src     string `taml:"role,omitempty"`
	Role    string `yaml:"role,omitempty"`
	Version string `yaml:"version,omitempty"`
	Name    string `yaml:"name,omitempty"`
}

type RoleListRequirement struct {
	Role map[string]Requirement
}

func (r *RoleRequirement) toRequirement() *Requirement {
	log.WithFields(log.Fields{
		"r": fmt.Sprintf("%+v", *r),
	}).Debug("Converting from RoleRequirement{}")

	req := Requirement{}

	if r.Src != "" {
		req.Src = r.Src
	} else if r.Role != "" {
		split := strings.Split(r.Role, ".")
		r.Src, _, _ = galaxy.FindRoleURL(split[0], split[1])
	} else if r.Name != "" {
		split := strings.Split(r.Name, ".")
		r.Src, _, _ = galaxy.FindRoleURL(split[0], split[1])
	}

	req.Src = strings.TrimSuffix(req.Src, ".git")
	req.Version = r.Version
	return &req
}

type SrcRequirement struct {
	Src     string `yaml:"src"`
	Version string `yaml:"version,omitempty"`
}

func (r *SrcRequirement) toRequirement() *Requirement {
	log.WithFields(log.Fields{
		"src": r.Src,
	}).Debug("Converting from SrcRequirement{}")

	req := Requirement{}

	if strings.HasPrefix(r.Src, "https://") {
		req.Src = strings.TrimSuffix(r.Src, ".git")
	} else {
		parts := strings.Split(r.Src, ".")
		log.WithFields(log.Fields{
			"parts": parts,
		}).Debug("Finding role url")

		req.Src, _, _ = galaxy.FindRoleURL(parts[0], parts[1])
	}

	req.Version = r.Version
	return &req
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

	iface := loadRequirementsYAML(requirementsData)

	for _, item := range iface {
		fmt.Println(fmt.Sprintf("%+v", item))
		itemJSON, err := json.MarshalIndent(item, "", "    ")
		if err != nil {
			panic(err)
		}

		req := convertAnythingToRequirement(itemJSON)
		*r = append(*r, *req)
	}
}

func convertAnythingToRequirement(in []byte) *Requirement {
	var roleReq RoleRequirement
	err := json.Unmarshal(in, &roleReq)
	if err == nil && (roleReq.Role != "" || roleReq.Name != "") {
		return roleReq.toRequirement()
	}

	var srcReq SrcRequirement
	err = json.Unmarshal(in, &srcReq)
	if err == nil {
		return srcReq.toRequirement()
	}

	return nil
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
