package docker

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

func Build(args []string) {
	command := fmt.Sprintf("%s %s --no-cache %s",
		args[0],
		args[1],
		strings.Join(args[2:], " "),
	)

	fmt.Println(fmt.Sprintf("command: %+v", command))

	fmt.Println(alpine(strings.Split(bash(command), "\n")))

}

func alpine(lines []string) string {
	var packages []string
	allPackages := map[string]string{}

	for _, line := range lines {
		parts := strings.Split(line, " ")

		packages = append(packages, extractAlpinePackages(parts)...)

		// (3/3) Installing htop (2.2.0-r0)
		if len(parts) < 2 || parts[1] != "Installing" {
			continue
		}

		allPackages[parts[2]] = sanitiseAlpineVersion(parts[3])
	}

	requested := map[string]string{}

	for _, pack := range packages {
		if value, ok := allPackages[pack]; ok == true {
			requested[pack] = value
		}

	}

	var ret []string
	for k, v := range requested {
		ret = append(ret, fmt.Sprintf("%s=%s", k, v))
	}

	sort.Strings(ret)
	return strings.Join(ret, "\n")
}

func extractAlpinePackages(parts []string) []string {
	apk := false
	add := false
	var ret []string

	for _, arg := range parts {
		if "apk" == arg {
			apk = true
			continue
		}
		if "add" == arg {
			add = true
			continue
		}
		if apk && add {
			ret = append(ret, arg)
		}
	}
	return ret
}

func sanitiseAlpineVersion(version string) string {
	version = strings.TrimLeft(version, "(")
	version = strings.TrimRight(version, ")")
	return version
}

func bash(command string) string {
	cmd := exec.Command("bash", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	outStr, _ := string(stdout.Bytes()), string(stderr.Bytes())

	return outStr
}
