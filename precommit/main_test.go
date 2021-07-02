package precommit

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func slurp(t *testing.T, file string) []byte {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	return bytes
}

func TestUpdate(t *testing.T) {
	var cases = []struct {
		name   string
		input  []byte
		output []byte
		err    error
	}{
		{
			name:  "valid precommit config file",
			input: slurp(t, "../.pre-commit-config.yaml"),
			err:   nil,
		},
		{
			name:   "valid precommit config file",
			input:  slurp(t, "../test/fixtures/pre-commit.yaml"),
			output: slurp(t, "../test/fixtures/pre-commit.updated.yaml"),
			err:    nil,
		},
		{
			name:   "invalid precommit file",
			input:  slurp(t, "../.github/workflows/pr.yml"),
			output: []byte{},
			err:    ErrorEmptyReposConfig,
		},
	}

	ghToken := os.Getenv("GITHUB_READONLY_TOKEN")
	if ghToken == "" {
		t.Fatal("Error. GITHUB_READONLY_TOKEN not set")
	}

	for _, test := range cases {
		output, err := Update(test.input, ghToken)
		assert.Equal(t, test.err, err, test.name)
		if test.output != nil {
			assert.Equal(t, test.output, []byte(output), test.name)
		}
	}
}
