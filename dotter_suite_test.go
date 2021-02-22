package dotter_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDotter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dotter Suite")
}

func mkdir(pathParts ...string) string {
	path := filepath.Join(pathParts...)
	os.MkdirAll(path, os.ModePerm)
	return path
}

func tmpdir() string {
	dir, _ := ioutil.TempDir("", "dotter")
	return dir
}

func writeFile(dir string, file string, content string) {
	ioutil.WriteFile(filepath.Join(dir, file), []byte(content), 0666)
}

func rm(dir string, file string) {
	os.Remove(filepath.Join(dir, file))
}

func rmdir(dir string) {
	os.RemoveAll(dir)
}
