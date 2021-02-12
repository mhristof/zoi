package gh

import (
	"errors"
	"regexp"
	"strings"

	"mvdan.cc/xurls/v2"
)

func Release(line string) string {
	var parsers = []func(string) (*Url, error){
		parseGit,
		parseHttp,
	}

	for _, parser := range parsers {
		gURL, err := parser(line)
		if err == nil {
			next, err := gURL.NextRelease()
			if err != nil {
				return line
			}

			return strings.Replace(line, gURL.Url, next, -1)
		}
	}

	return line
}

func parseGit(line string) (*Url, error) {
	regex := `git@github.com.*ref=[\w\.]*`
	re := regexp.MustCompile(regex)
	found := re.Find([]byte(line))

	if len(found) == 0 {
		return nil, errors.New("Not a git@github.com url")
	}

	gURL, err := ParseUrl(string(found))
	if err != nil {
		return nil, err
	}

	return gURL, nil
}

func parseHttp(line string) (*Url, error) {
	urls := xurls.Relaxed()
	url := urls.FindString(line)

	if url == "" {
		return nil, errors.New("No URL found")
	}

	if strings.Contains(url, "releases/latest") {
		return nil, errors.New("Not a release url")
	}

	gURL, err := ParseUrl(url)
	if err != nil {
		return nil, err
	}

	return gURL, nil
}
