package main

import (
	"testing"

	"github.com/mhristof/zoi/ansible"
	"github.com/mhristof/zoi/github"
	"github.com/stretchr/testify/assert"
)

func TestLatestTag(t *testing.T) {
	var cases = []struct {
		name     string
		in       ansible.Requirement
		expected ansible.Requirement
	}{
		{
			"unset version",
			ansible.Requirement{
				Src:     "https://github.com/geerlingguy/ansible-role-jenkins",
				Version: "",
			},
			ansible.Requirement{
				Src:     "https://github.com/geerlingguy/ansible-role-jenkins",
				Version: "4.2.1",
			},
		},
		{
			"outdated version",
			ansible.Requirement{
				Src:     "https://github.com/geerlingguy/ansible-role-jenkins",
				Version: "4.0.0",
			},
			ansible.Requirement{
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
