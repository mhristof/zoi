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
		name  string
		input []byte
		err   error
	}{
		{
			name:  "valid precommit config file",
			input: slurp(t, "../.pre-commit-config.yaml"),
			err:   nil,
		},
	}

	ghToken := os.Getenv("GITHUB_READONLY_TOKEN")
	if ghToken == "" {
		t.Fatal("Error. GITHUB_READONLY_TOKEN not set")
	}

	for _, test := range cases {
		_, err := Update(test.input, ghToken)
		assert.Equal(t, test.err, err, test.name)
	}
}
