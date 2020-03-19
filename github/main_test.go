package github

//
// requirements_test.go
// Copyright (C) 2020 mhristof <mhristof@Mikes-MBP>
//
// Distributed under terms of the MIT license.
//

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractUserRepoFromSrc(t *testing.T) {
	var cases = []struct {
		src  string
		user string
		repo string
	}{
		{
			src:  "geerlingguy.jenkins",
			user: "geerlingguy",
			repo: "ansible-role-jenkins",
		},
		{
			src:  "https://github.com/geerlingguy/ansible-role-jenkins",
			user: "geerlingguy",
			repo: "ansible-role-jenkins",
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.src, func(t *testing.T) {
			t.Parallel() // marks each test case as capable of running in parallel with each other
			t.Log(tt.src)
			rUser, rRepo := extractUserRepoFromSrc(tt.src)
			assert.Equal(t, rUser, tt.user)
			assert.Equal(t, rRepo, tt.repo)
		})
	}
}

func TestLatestTag(t *testing.T) {
	var cases = []struct {
		src string
		tag string
	}{
		{
			src: "https://github.com/mhristof/cautious-potato",
			tag: "1.2",
		},
	}

	g := New()

	for _, tt := range cases {
		tt := tt
		t.Run(tt.src, func(t *testing.T) {
			t.Parallel() // marks each test case as capable of running in parallel with each other
			t.Log(tt.src)
			assert.Equal(t, g.LatestTag(tt.src), tt.tag)
		})
	}

}

func TestLatestBranchCommit(t *testing.T) {
	var cases = []struct {
		repo   string
		commit string
	}{
		{
			repo:   "https://github.com/mhristof/cautious-potato",
			commit: "1e5116344fd677a646e1c559a832bb47786c9c99",
		},
	}

	g := New()

	for _, tt := range cases {
		tt := tt
		t.Run(tt.repo, func(t *testing.T) {
			t.Parallel() // marks each test case as capable of running in parallel with each other
			t.Log(tt.repo)
			assert.Equal(t, g.LatestBranchCommit(tt.repo), tt.commit)
		})
	}
}
