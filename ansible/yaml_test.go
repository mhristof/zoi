// +build !fast

package ansible

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"testing"
)

func TestYAML(t *testing.T) {
	folder := "requirement-yamls"
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".yml") {
			continue
		}
		path := filepath.Join(folder, file.Name())
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			t.Log(path)
			reqs := Requirements{}
			reqs.LoadFromFile(path)
			reqs.Update().SaveToFile("/dev/null")
		})
	}

}
