package docker

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
)

func toFile(contents string) string {
	// create a new dir so that the whole universe is not uploaded when we
	// docker build.
	dir, err := ioutil.TempDir("", "zoi")
	if err != nil {
		log.Fatal(err)
	}

	file, err := ioutil.TempFile(dir, "Dockerfile.")
	if err != nil {
		panic(err)
	}

	_, err = file.Write([]byte(contents))
	if err != nil {
		panic(err)
	}

	return file.Name()
}

func TestNew(t *testing.T) {
	var cases = []struct {
		name     string
		content  string
		expected string
	}{
		{
			name: "simple docker file with one apk package",
			content: heredoc.Doc(`
				FROM alpine:3.12.0

				RUN apk add bash
			`),
			expected: heredoc.Doc(`
				FROM alpine:3.12.0

				RUN apk add bash=5.0.17-r0
			`),
		},
		{
			name: "a few apk packages",
			content: heredoc.Doc(`
				FROM alpine:3.12.0

				RUN apk add bash jq python3
			`),
			expected: heredoc.Doc(`
				FROM alpine:3.12.0

				RUN apk add bash=5.0.17-r0 jq=1.6-r1 python3=3.8.3-r0
			`),
		},
		{
			name: "user pinned package",
			content: heredoc.Doc(`
				FROM alpine:3.12.0

				RUN apk add -f bash jq python3=3.8.2-r0
			`),
			expected: heredoc.Doc(`
				FROM alpine:3.12.0

				RUN apk add -f bash=5.0.17-r0 jq=1.6-r1 python3=3.8.2-r0
			`),
		},
		{
			name: "multi line RUN",
			content: heredoc.Doc(`
				FROM alpine:3.12.0

				RUN apk add -f bash jq python3 &&\
					date
			`),
			expected: heredoc.Doc(`
				FROM alpine:3.12.0

				RUN apk add -f bash=5.0.17-r0 jq=1.6-r1 python3=3.8.3-r0 &&\
					date
			`),
		},
	}

	for _, test := range cases {
		f := New(toFile(test.content))
		assert.Equal(t, test.expected, f.Render(), test.name)
	}
}
