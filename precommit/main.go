package precommit

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mhristof/zoi/gh"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Repo struct {
	Repo  string `yaml:"repo,omitempty"`
	Rev   string `yaml:"rev,omitempty"`
	Hooks []struct {
		Args    []string `yaml:"args,omitempty"`
		ID      string   `yaml:"id,omitempty"`
		Name    string   `yaml:"name,omitempty"`
		Exclude string   `yaml:"exclude,omitempty"`
	} `yaml:"hooks,omitempty"`
}

type Config struct {
	Repos []*Repo `yaml:"repos,omitempty"`
}

func Update(bytesIn []byte, token string) (string, error) {
	var config Config

	err := yaml.Unmarshal(bytesIn, &config)
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

	b := new(bytes.Buffer)
	yamlEncoder := yaml.NewEncoder(b)
	yamlEncoder.SetIndent(2)

	err = yamlEncoder.Encode(&config)
	if err != nil {
		return "", errors.Wrap(err, "Cannot encode config")
	}

	return strings.Join([]string{"---", b.String()}, "\n"), nil
}
