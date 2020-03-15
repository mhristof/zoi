//
// main.go
// Copyright (C) 2020 mhristof <mhristof@Mikes-MBP>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"fmt"
	"os"

	"github.com/mhristof/zoi/ansible"
	"github.com/mhristof/zoi/github"
)

func main() {
	if len(os.Args) < 2 {
		panic("Error, expected one argument")
	}

	requirementsPath := os.Args[1]
	if _, err := os.Stat(requirementsPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("Error, file %s does not exist", requirementsPath))
	}

	reqs := ansible.Requirements{}
	reqs.LoadFromFile(requirementsPath)

	gh := github.New()

	var latestRequirements ansible.Requirements
	for _, requirement := range reqs {
		latest := gh.LatestTag(requirement.Src)
		latestRequirements = append(latestRequirements, ansible.Requirement{
			Src:     requirement.Src,
			Version: latest,
		})
	}

	latestRequirements.SaveToFile("latest.yml")
}
