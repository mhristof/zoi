package gh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			out:  "lorem ipsum https://github.com/mhristof/semver/releases/download/v0.5.0/semver.darwin test",
		},
		{
			name: "latest release url",
			in:   "https://github.com/mhristof/checkov2vim/releases/latest/download/checkov2vim",
			out:  "https://github.com/mhristof/checkov2vim/releases/latest/download/checkov2vim",
		},
	}

	for _, test := range cases {
		assert.Equal(t, test.out, Release(test.in), test.name)

	}
}
