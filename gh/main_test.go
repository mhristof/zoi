package gh

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var semverLatest = "0.5.0"

func TestRelease(t *testing.T) {
	var cases = []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "string without url",
			in:   "this is a test",
			out:  "this is a test",
		},
		{
			name: "github url",
			in:   "lorem ipsum https://github.com/mhristof/semver/releases/download/v0.3.2/semver.darwin test",
			out:  fmt.Sprintf("lorem ipsum https://github.com/mhristof/semver/releases/download/v%s/semver.darwin test", semverLatest),
		},
		{
			name: "latest release url",
			in:   "https://github.com/mhristof/checkov2vim/releases/latest/download/checkov2vim",
			out:  "https://github.com/mhristof/checkov2vim/releases/latest/download/checkov2vim",
		},
		{
			name: "github ssh url with ?ref=",
			in:   "git@github.com:mhristof/semver.git?ref=v0.3.2",
			out:  fmt.Sprintf("git@github.com:mhristof/semver.git?ref=v%s", semverLatest),
		},
	}

	for _, test := range cases {
		assert.Equal(t, test.out, Release(test.in), test.name)

	}
}
