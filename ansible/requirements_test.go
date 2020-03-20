package ansible

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/mhristof/zoi/github"
	"github.com/stretchr/testify/assert"
)

func TestLatestTag(t *testing.T) {
	var cases = []struct {
		name     string
		in       Requirement
		expected Requirement
	}{
		{
			"unset version",
			Requirement{
				Src:     "https://github.com/geerlingguy/ansible-role-jenkins",
				Version: "",
			},
			Requirement{
				Src:     "https://github.com/geerlingguy/ansible-role-jenkins",
				Version: "4.2.1",
			},
		},
		{
			"outdated version",
			Requirement{
				Src:     "https://github.com/geerlingguy/ansible-role-jenkins",
				Version: "4.0.0",
			},
			Requirement{
				Src:     "https://github.com/geerlingguy/ansible-role-jenkins",
				Version: "4.2.1",
			},
		},
	}

	g := github.New()

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // marks each test case as capable of running in parallel with each other
			t.Log(tt.in)
			assert.Equal(t, g.LatestTag(tt.in.Src), tt.expected.Version, "they should be equal")
		})
	}
}

func TestUpdate(t *testing.T) {
	var cases = []struct {
		name string
		yaml string
		res  Requirement
	}{
		{
			name: "role with valid github src, name, scm and version master",
			yaml: heredoc.Doc(`
			- src: 'https://github.com/mhristof/cautious-potato'
			  name: 'mhristof.cautious-potato'
			  scm: 'git'
			`),
			res: Requirement{
				Src:     "https://github.com/mhristof/cautious-potato",
				Version: "1.2",
			},
		},
		{
			name: "role with role",
			yaml: heredoc.Doc(`
			- role: snakeego.docker
			`),
			res: Requirement{
				Src:     "https://github.com/snakeego/ansible-role-docker",
				Version: "1.3.0",
			},
		},
		{
			name: "role with name",
			yaml: heredoc.Doc(`
			- name: snakeego.docker
			`),
			res: Requirement{
				Src:     "https://github.com/snakeego/ansible-role-docker",
				Version: "1.3.0",
			},
		},
		{
			name: "role with git+ in the src",
			yaml: heredoc.Doc(`
			- src: git+https://github.com/danie1cohen/ansible-virtualenv3
			`),
			res: Requirement{
				Src:     "https://github.com/danie1cohen/ansible-virtualenv3",
				Version: "05488949b99bd74d53b77b086a32572d9af0eaeb",
			},
		},
		{
			name: "roles dictionary with a role that has name defined",
			yaml: heredoc.Doc(`
			roles:
			  - name: snakeego.docker
			`),
			res: Requirement{
				Src:     "https://github.com/snakeego/ansible-role-docker",
				Version: "1.3.0",
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			t.Log(tt.name)
			var r Requirements
			r.LoadBytes([]byte(tt.yaml))
			r = *r.Update()
			assert.Equal(t, tt.res, r[0])
		})
	}

}
