package ansible

import (
	"io/ioutil"

	"github.com/mhristof/zoi/github"
	"github.com/mhristof/zoi/log"
	"gopkg.in/yaml.v3"
)

type Requirement struct {
	Src     string `yaml:"src,omitempty"`
	Version string `yaml:"version,omitempty"`
}

type Requirements []Requirement

func (r *Requirements) LoadFromFile(path string) {
	requirementsData, err := ioutil.ReadFile(path)
	err = yaml.Unmarshal(requirementsData, r)
	if err != nil {
		log.WithFields(log.Fields{
			"path": path,
			"err":  err,
		}).Error("Error while loading yaml file")
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
		latestRequirements = append(latestRequirements, Requirement{
			Src:     requirement.Src,
			Version: latest,
		})
	}
	return &latestRequirements
}
