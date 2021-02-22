package step_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/fatih/color"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	color.NoColor = true
}

func TestStep(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Step Suite")
}

func mkdir(pathParts ...string) string {
	path := filepath.Join(pathParts...)
	os.MkdirAll(path, os.ModePerm)
	return path
}

func ln(linkPath string, target string) {
	os.Symlink(target, linkPath)
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

type TestExecutor struct {
	target   string
	source   string
	backedUp []string
	infoLog  []string
	errorLog []string
}

func NewTestExecutor(source string, target string) *TestExecutor {
	return &TestExecutor{
		target:   target,
		source:   source,
		backedUp: make([]string, 0),
		infoLog:  make([]string, 0),
		errorLog: make([]string, 0),
	}
}

func (exec TestExecutor) GetTargetPath(path string) string {
	return filepath.Join(exec.target, path)
}

func (exec TestExecutor) GetSourcePath(path string) string {
	return filepath.Join(exec.source, path)
}

func (exec *TestExecutor) ForceRemove(path string) error {
	exec.backedUp = append(exec.backedUp, path)
	return os.RemoveAll(path)
}

func (exec *TestExecutor) PrintInfo(message string) {
	exec.infoLog = append(exec.infoLog, message)
}

func (exec *TestExecutor) PrintError(message string) {
	exec.errorLog = append(exec.errorLog, message)
}
