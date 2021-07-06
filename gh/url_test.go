package gh

import (
	"fmt"
	"os"
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
			err:  ErrorURLTooShort,
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
			name: "next release for mhristof/semver repo",
			repo: Url{
				Host:    "https://github.com",
				Owner:   "mhristof",
				Repo:    "semver",
				Release: "v0.1.0",
				Url:     "v0.1.0",
			},
			out: fmt.Sprintf("v%s", semverLatest),
		},
		{
			name: "binary containing the version",
			repo: Url{
				Host:    "https://github.com",
				Owner:   "mhristof",
				Repo:    "zoi-cli-cli",
				Release: "v1.7.0",
				Url:     "https://github.com/mhristof/zoi-cli-cli/releases/download/v1.7.0/gh_1.7.0_$(GH_OS)_amd64.tar.gz",
			},
			out: "https://github.com/mhristof/zoi-cli-cli/releases/download/v1.9.1/gh_1.9.1_$(GH_OS)_amd64.tar.gz",
		},
		{
			name: "github actions url",
			repo: Url{
				Host:    "https://github.com",
				Owner:   "mhristof",
				Repo:    "zoi-github-autopr",
				Release: "0.1.1",
				Url:     "mhristof/zoi-github-autopr@0.1.1",
			},
			out: "mhristof/zoi-github-autopr@0.2.0",
		},
		{
			name: "short version that is contained in the new version string as well",
			repo: Url{
				Host:    "https://github.com",
				Owner:   "mhristof",
				Repo:    "zoi-checkout",
				Release: "v2",
				Url:     "mhristof/zoi-checkout@v2",
			},
			out: "mhristof/zoi-checkout@v2.3.4",
		},
		{
			name: "project with an older release and a newer tag",
			repo: Url{
				Host:    "https://github.com",
				Owner:   "mhristof",
				Repo:    "zoi-go-humanize",
				Release: "v1.0.0",
				Url:     "mhristof/zoi-go-humanize?ref=v1.0.0",
			},
			out: "mhristof/zoi-go-humanize?ref=v1.0.0",
		},
	}

	ghToken := os.Getenv("GITHUB_READONLY_TOKEN")
	if ghToken == "" {
		t.Fatal("Error. GITHUB_READONLY_TOKEN not set")
	}

	for _, test := range cases {
		test.repo.Token = ghToken
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
