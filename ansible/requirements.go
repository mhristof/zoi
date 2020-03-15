package ansible

import (
	"io/ioutil"

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
