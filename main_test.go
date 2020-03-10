package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

	var g GitHub
	g.New()

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // marks each test case as capable of running in parallel with each other
			t.Log(tt.in)
			assert.Equal(t, g.latestTag(tt.in), &tt.expected, "they should be equal")
		})
	}
}
