package dotter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var (
	cYellow  = color.New(color.FgYellow).SprintFunc()
	cbYellow = color.New(color.FgYellow, color.Bold).SprintFunc()
	cGreen   = color.New(color.FgGreen).SprintFunc()
	cbGreen  = color.New(color.FgGreen, color.Bold).SprintFunc()
	cbRed    = color.New(color.FgRed, color.Bold).SprintFunc()
	cbCyan   = color.New(color.FgCyan, color.Bold).SprintFunc()
)

type Executor struct {
	SourceDirectory string
	TargetDirectory string
	Configuration   Configuration
}

func NewExecutor(sourceDirectory string, targetDirectory string, config Configuration) Executor {
	exec := Executor{}
	exec.SourceDirectory = sourceDirectory
	exec.TargetDirectory = targetDirectory
	exec.Configuration = config
	return exec
}

func (exec Executor) output(format string, a ...interface{}) {
	if !exec.Configuration.Options.Quiet {
		fmt.Printf(format, a...)
	}
}

func (exec Executor) Execute() error {
	exec.output(
		cYellow("Installing %s to %s\n"),
		cbYellow(exec.SourceDirectory),
		cbYellow(exec.TargetDirectory),
	)
	if exec.Configuration.SourcePath != "" {
		exec.output(
			cYellow("Using %s\n"),
			cbYellow(exec.Configuration.SourcePath),
		)
	}

	for _, step := range exec.Configuration.Steps {
		exec.output(
			cGreen("%s: %s\n"),
			step.GetActivityLabel(),
			cbGreen(step.GetActivityDetails()),
		)
		err := step.Execute(exec)

		if err != nil {
			// TODO print error
			if exec.Configuration.Options.StopOnError {
				return err
			}
		}
	}

	exec.output(cbYellow("Complete.\n"))
	return nil
}

func (exec Executor) GetTargetPath(path string) string {
	return filepath.Join(exec.TargetDirectory, path)
}

func (exec Executor) GetSourcePath(path string) string {
	return filepath.Join(exec.SourceDirectory, path)
}

func (exec Executor) ForceRemove(path string) error {
	if exec.Configuration.Options.BackupForced != "" {
		fmt.Printf("Backing up %s...\n", path)
	}

	return os.RemoveAll(path)
}

func (exec Executor) PrintInfo(message string) {
	for _, line := range indentString(message) {
		exec.output("%s\n", cbCyan(line))
	}
}

func (exec Executor) PrintError(message string) {
	for _, line := range indentString(message) {
		exec.output("%s\n", cbRed(line))
	}
}

func indentString(value string) []string {
	lines := strings.Split(value, "\n")
	indented := make([]string, 0, len(lines))

	for idx, line := range lines {
		if line == "" && idx == (len(lines)-1) {
			break
		}
		indented = append(indented, "    "+line)
	}

	return indented
}
