package gh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAction(t *testing.T) {
	var cases = []struct {
		name   string
		input  string
		output *Url
		err    bool
	}{
		{
			name:  "simple actions repository",
			input: "jessfraz/branch-cleanup-action@master",
			output: &Url{
				Host:    "https://github.com",
				Owner:   "jessfraz",
				Repo:    "branch-cleanup-action",
				Release: "master",
				Url:     "jessfraz/branch-cleanup-action@master",
			},
		},
		{
			name:  "non actions url",
			input: "git@github.com:mhristof/zoi.git",
			err:   true,
		},
	}

	for _, test := range cases {
		res, err := parseAction(test.input)
		assert.Equal(t, test.err, err != nil, test.name)
		assert.Equal(t, test.output, res, test.name)
	}
}
