package precommit

import (
	"fmt"
	"strings"

	"github.com/mhristof/zoi/gh"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Repo struct {
	Hooks []struct {
		Args []string `yaml:"args,omitempty"`
		ID   string   `yaml:"id,omitempty"`
		Name string   `yaml:"name,omitempty"`
	} `yaml:"hooks,omitempty"`
	Repo string `yaml:"repo,omitempty"`
	Rev  string `yaml:"rev,omitempty"`
}

type Config struct {
	Repos []*Repo `yaml:"repos,omitempty"`
}

func Update(bytes []byte, token string) (string, error) {
	var config Config

	err := yaml.Unmarshal(bytes, &config)
	if err != nil {
		return "", errors.Wrap(err, "Cannot unmarshal config")
	}

	for _, repo := range config.Repos {
		latest := strings.TrimPrefix(
			gh.Release(fmt.Sprintf("%s?ref=%s", repo.Repo, repo.Rev), token),
			fmt.Sprintf("%s?ref=", repo.Repo),
		)
		repo.Rev = latest
	}

	out, err := yaml.Marshal(config)
	if err != nil {
		return "", errors.Wrap(err, "Cannot marshal yaml")
	}

	return string(out), nil
}
