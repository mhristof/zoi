package gh

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHttpUrl(t *testing.T) {
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
		{
			name: "https url with ref=",
			in:   "https://github.com/mhristof/terraform-aws-vpc-1?ref=v0.1.2",
			out: &Url{
				Host:    "https://github.com",
				Owner:   "mhristof",
				Repo:    "terraform-aws-vpc-1",
				Url:     "https://github.com/mhristof/terraform-aws-vpc-1?ref=v0.1.2",
				Release: "v0.1.2",
			},
		},
	}

	for _, test := range cases {
		url, err := ParseHttpUrl(test.in)
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
				Host:    "github.com",
				Owner:   "mhristof",
				Repo:    "semver",
				Release: "v1.2.3",
				Url:     "git@github.com:mhristof/semver.git?ref=v1.2.3",
			},
		},
		// {
		// 	name: "invalid github ssh url",
		// 	in:   "git@gitlab.com:mhristof/semver.git?ref=v1.2.3",
		// 	out:  nil,
		// 	err:  ErrorWrongHost,
		// },
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
			out: fmt.Sprintf("v%s", semverLatest),
		},
	}

	for _, test := range cases {
		next, _ := test.repo.NextRelease()
		assert.Equal(t, test.out, next, test.name)

	}
}

func Test(t *testing.T) {
	var cases = []struct {
		name string
		in   []string
		out  string
	}{
		{
			name: "release download url",
			in:   strings.Split("https://github.com/mhristof/semver/releases/download/v0.5.0/semver.darwin", "/"),
			out:  "v0.5.0",
		},
		{
			name: "http url with ref",
			in:   strings.Split("https://github.com/mhristof/terraform-aws-vpc-1?ref=v0.1.2", "/"),
			out:  "v0.1.2",
		},
	}

	for _, test := range cases {
		assert.Equal(t, test.out, getRelease(test.in...), test.name)
	}
}
