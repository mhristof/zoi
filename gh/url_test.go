package gh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUrl(t *testing.T) {
	var cases = []struct {
		name string
		in   string
		out  *Url
		err  error
	}{
		{
			name: "valid url",
			in:   "https://github.com/mhristof/semver",
			out: &Url{
				Host:  "https://github.com",
				Owner: "mhristof",
				Repo:  "semver",
				Url:   "https://github.com/mhristof/semver",
			},
		},
		{
			name: "valid url with release",
			in:   "https://github.com/mhristof/semver/releases/download/v1.2.3",
			out: &Url{
				Host:    "https://github.com",
				Owner:   "mhristof",
				Repo:    "semver",
				Release: "v1.2.3",
				Url:     "https://github.com/mhristof/semver/releases/download/v1.2.3",
			},
		},
		{
			name: "short github url",
			in:   "https://github.com/mhristof",
			out:  nil,
			err:  ErrorUrlTooShort,
		},
		{
			name: "wrong host",
			in:   "https://githu.com/mhristof",
			out:  nil,
			err:  ErrorWrongHost,
		},
		{
			name: "valid url with .git in the project name",
			in:   "https://github.com/VundleVim/Vundle.vim.git",
			out: &Url{
				Host:  "https://github.com",
				Owner: "VundleVim",
				Repo:  "Vundle.vim",
				Url:   "https://github.com/VundleVim/Vundle.vim.git",
			},
		},
	}

	for _, test := range cases {
		url, err := ParseUrl(test.in)
		assert.Equal(t, err, test.err, test.name)
		assert.Equal(t, test.out, url, test.name)

	}
}

func TestParseGitUrl(t *testing.T) {
	var cases = []struct {
		name string
		in   string
		out  *Url
		err  error
	}{
		{
			name: "valid github ssh url with ref",
			in:   "git@github.com:mhristof/semver.git?ref=v1.2.3",
			out: &Url{
				Host:    "git@github.com",
				Owner:   "mhristof",
				Repo:    "semver",
				Release: "v1.2.3",
				Url:     "git@github.com:mhristof/semver.git?ref=v1.2.3",
			},
		},
		{
			name: "invalid github ssh url",
			in:   "git@gitlab.com:mhristof/semver.git?ref=v1.2.3",
			out:  nil,
			err:  ErrorWrongHost,
		},
	}

	for _, test := range cases {
		url, err := ParseGitUrl(test.in)
		assert.Equal(t, err, test.err, test.name)
		assert.Equal(t, test.out, url, test.name)

	}
}

func TestNextReleaseUrl(t *testing.T) {
	var cases = []struct {
		name string
		repo Url
		out  string
	}{
		{
			name: "next remease for mhristof/semver repo",
			repo: Url{
				Host:    "https://github.com",
				Owner:   "mhristof",
				Repo:    "semver",
				Release: "v0.1.0",
				Url:     "v0.1.0",
			},
			out: "v0.5.0",
		},
	}

	for _, test := range cases {
		next, _ := test.repo.NextRelease()
		assert.Equal(t, test.out, next, test.name)

	}
}
