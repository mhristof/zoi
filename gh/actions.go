package gh

import (
	"errors"
	"regexp"
	"strings"
)

func parseAction(line string) (*Url, error) {
	// for example:
	//	jessfraz/branch-cleanup-action@master
	regex := `[\w-]+/[\w-]+@.*`

	re := regexp.MustCompile(regex)
	found := re.Find([]byte(line))

	if len(found) == 0 {
		return nil, errors.New("Not a github actions url")
	}

	fields := strings.FieldsFunc(string(found), func(r rune) bool {
		return r == '/' || r == '@'
	})
	return &Url{
		Host:    "https://github.com",
		Owner:   fields[0],
		Repo:    fields[1],
		Release: fields[2],
		Url:     string(found),
	}, nil
}
