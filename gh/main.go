package gh

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"mvdan.cc/xurls/v2"
)

func Release(line string) string {

	httpUrl, err := releaseHttp(line)

	if err == nil {
		return httpUrl
	}

	gitUrl, err := releaseGit(line)
	if err == nil {
		return gitUrl
	}

	return line
}

func releaseGit(line string) (string, error) {
	regex := `git@github.com.*ref=\w*`
	re := regexp.MustCompile(regex)
	found := re.Find([]byte(line))

	if len(found) == 0 {
		return "", errors.New("Not a git@github.com url")
	}

	gURL, err := ParseUrl(string(found))
	if err != nil {
		return "", err
	}

	fmt.Println(fmt.Sprintf("gURL: %+v", gURL))

	return "", nil
}

func releaseHttp(line string) (string, error) {
	urls := xurls.Relaxed()
	url := urls.FindString(line)

	if url == "" {
		return line, errors.New("No URL found")
	}

	if strings.Contains(url, "releases/latest") {
		return line, errors.New("Not a release url")
	}

	gUrl, err := ParseUrl(url)
	if err != nil {
		return line, err
	}

	next, err := gUrl.NextRelease()
	if err != nil {
		return line, err
	}

	return strings.Replace(line, url, next, -1), nil
}
