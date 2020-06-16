package docker

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mhristof/zoi/log"
)

type File struct {
	Path    string
	Changes map[string]string
}

type apkPackage struct {
	name    string
	version string
}

func New(file string) File {
	var ret = File{
		Path:    file,
		Changes: map[string]string{},
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.WithFields(log.Fields{
			"file": file,
		}).Error("File not found")
	}

	cmd := exec.Command("docker", "build", "--no-cache", "-t",
		fmt.Sprintf("zoi-%s", name(ret.Path)),
		"-f", ret.Path, filepath.Dir(ret.Path),
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.WithFields(log.Fields{
			"cmd":        cmd.Args,
			"err":        err,
			"cmd.Stderr": cmd.Stderr,
		}).Error("Cannot build container")

	}

	var step string
	for _, line := range strings.Split(string(stdout.Bytes()), "\n") {
		if strings.HasPrefix(line, "Step") {
			step = strings.TrimSpace(strings.Join(strings.Split(line, ":")[1:], " "))
			ret.Changes[step] = step
			continue
		}

		pack, err := apk(line)
		if err != nil {
			continue
		}

		if strings.Index(step, pack.name) == -1 {
			continue
		}

		ret.Changes[step] = strings.Replace(
			ret.Changes[step], pack.name, fmt.Sprintf("%s=%s", pack.name, pack.version), -1,
		)
	}

	return ret
}

func (f *File) Render() string {
	lines, err := ioutil.ReadFile(f.Path)
	if err != nil {
		panic(err)
	}

	var ret []string
	//for k, v := range f.Changes {
	//fmt.Println(fmt.Sprintf("changes[%s] = %s", k, v))
	//}
	for _, line := range strings.Split(strings.ReplaceAll(string(lines), "&&\\\n", `&&`), "\n") {
		if value, ok := f.Changes[line]; ok == true {
			ret = append(ret, strings.ReplaceAll(value, `&&`, "&&\\\n"))
			continue
		}
		ret = append(ret, line)
	}

	return strings.Join(ret, "\n")
}

func apk(line string) (apkPackage, error) {
	var pkg apkPackage
	reg := regexp.MustCompile(`^\(\d*/\d*\) Installing (.*) \((.*)\)$`)

	finds := reg.FindAllStringSubmatch(line, -1)
	if len(finds) != 1 {
		return apkPackage{}, errors.New("Not an apk install command")
	}

	pkg.name = finds[0][1]
	pkg.version = finds[0][2]

	return pkg, nil
}

func name(file string) string {
	f, err := os.Open(file)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"file": file,
		}).Error("Cannot read file")
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot io.Copy")

	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
