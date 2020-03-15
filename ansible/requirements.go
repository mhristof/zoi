package ansible

import (
	"io/ioutil"

	"github.com/mhristof/zoi/github"
	"gopkg.in/yaml.v2"
)

type Requirement struct {
	Src     string `yaml:"src"`
	Version string `yaml:"version"`
}

type Requirements []Requirement

func (r *Requirements) LoadFromFile(path string) {
	requirementsData, err := ioutil.ReadFile(path)
	err = yaml.Unmarshal(requirementsData, r)
	if err != nil {
		panic(err)
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
		latest := gh.LatestTag(requirement.Src)
		latestRequirements = append(latestRequirements, Requirement{
			Src:     requirement.Src,
			Version: latest,
		})
	}
	return &latestRequirements
}
