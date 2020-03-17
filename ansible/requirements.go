package ansible

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mhristof/zoi/github"
	"github.com/mhristof/zoi/log"
	"gopkg.in/yaml.v3"
)

// type GenericReqiurement struct {
// 	data [string]interface{}
// }

// type GenericRequirements []GenericReqiurement

type Requirement struct {
	Src     string `yaml:"src,omitempty"`
	Version string `yaml:"version,omitempty"`
}

type Requirements []Requirement

type RoleRequirement struct {
	Role    string `yaml:"role"`
	Version string `yaml:"version,omitempty"`
}

func (r *RoleRequirement) toRequirement() *Requirement {
	log.WithFields(log.Fields{
		"r": fmt.Sprintf("%+v", *r),
	}).Debug("Converting to Requirement{}")

	req := Requirement{}

	parts := strings.Split(r.Role, ".")
	req.Src = fmt.Sprintf("https://github.com/%s/ansible-role-%s", parts[0], parts[1])
	req.Version = r.Version
	return &req
}

type SrcRequirement struct {
	Src     string `yaml:"src"`
	Version string `yaml:"version,omitempty"`
}

func (r *SrcRequirement) toRequirement() *Requirement {
	req := Requirement{}

	parts := strings.Split(r.Src, ".")
	req.Src = fmt.Sprintf("https://github.com/%s/ansible-role-%s", parts[0], parts[1])
	req.Version = r.Version
	return &req
}

func (r *Requirements) LoadFromFile(path string) {
	log.WithFields(log.Fields{
		"path": path,
	}).Debug("Loading requirements file")

	requirementsData, err := ioutil.ReadFile(path)

	var iface []interface{}
	err = yaml.Unmarshal(requirementsData, &iface)
	if err != nil {
		log.WithFields(log.Fields{
			"path": path,
			"err":  err,
		}).Error("Error while loading yaml file")
	}

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
	// dont forget to import "encoding/json"
	if err == nil && roleReq.Role != "" {
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
		latestRequirements = append(latestRequirements, Requirement{
			Src:     requirement.Src,
			Version: latest,
		})
	}
	return &latestRequirements
}
