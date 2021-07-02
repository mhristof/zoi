package precommit

import (
	"bytes"
	"fmt"
	"strings"

	pr "github.com/mhristof/go-precommit"
	"github.com/mhristof/zoi/gh"
	"github.com/mhristof/zoi/log"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Repos []*pr.Repo `yaml:"repos,omitempty"`
}

var (
	ErrorEmptyReposConfig = errors.New("Empty `repos` field")
)

func Update(bytesIn []byte, token string) (string, error) {
	var config Config

	err := yaml.Unmarshal(bytesIn, &config)
	if err != nil {
		return "", errors.Wrap(err, "Cannot unmarshal config")
	}

	if len(config.Repos) == 0 {
		return "", ErrorEmptyReposConfig
	}

	log.WithFields(log.Fields{
		"err": err,
	}).Debug("Handling a precommit file")

	for _, repo := range config.Repos {
		latest := strings.TrimPrefix(
			gh.Release(fmt.Sprintf("%s?ref=%s", repo.Repo, repo.Rev), token),
			fmt.Sprintf("%s?ref=", repo.Repo),
		)
		repo.Rev = latest
	}

	b := new(bytes.Buffer)
	yamlEncoder := yaml.NewEncoder(b)
	yamlEncoder.SetIndent(2)

	err = yamlEncoder.Encode(&config)
	if err != nil {
		return "", errors.Wrap(err, "Cannot encode config")
	}

	return strings.Join([]string{"---", b.String()}, "\n"), nil
}
