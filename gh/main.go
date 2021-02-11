package gh

import (
	"strings"

	"mvdan.cc/xurls/v2"
)

func Release(line string) string {
	urls := xurls.Relaxed()
	url := urls.FindString(line)

	if url == "" {
		return line
	}

	if strings.Contains(url, "releases/latest") {
		return line
	}

	gUrl, err := ParseUrl(url)
	if err != nil {
		return line
	}

	next, err := gUrl.NextRelease()
	if err != nil {
		return line
	}

	return strings.Replace(line, url, next, -1)
}
