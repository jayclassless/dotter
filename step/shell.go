package step

import (
	"bytes"
	"os"
	osexec "os/exec"
)

// ShellOptions contains non-command options for Shell steps
type ShellOptions struct {
	Quiet bool
}

// NewShellOptions creates a new instance of a ShellOptions struct
func NewShellOptions() ShellOptions {
	opt := ShellOptions{}
	opt.Quiet = true
	return opt
}

// ShellStep contains the specification for Shell steps
type ShellStep struct {
	ShellOptions `yaml:",inline"`
	Command      string
	Description  string
}

// NewShellStep creates a new instance of a ShellStep struct using default options
func NewShellStep() ShellStep {
	return NewShellStepWithDefaults(NewShellOptions())
}

// NewShellStepWithDefaults creates a new instance of a ShellStep struct using the specified options
func NewShellStepWithDefaults(defaults ShellOptions) ShellStep {
	step := ShellStep{}
	step.ShellOptions = defaults
	return step
}

// GetActivityLabel returns a short description of what a ShellStep does
func (step ShellStep) GetActivityLabel() string {
	return "Executing"
}

// GetActivityDetails returns description specific to this particular instance of the ShellStep
func (step ShellStep) GetActivityDetails() string {
	if step.Description != "" {
		return step.Description
	}
	return step.Command
}

// Execute runs the specified command in a shell
func (step ShellStep) Execute(exec StepExecutor) error {
	cmd := osexec.Command(
		step.getShell(),
		"-c",
		step.Command,
	)

	cmd.Dir = exec.GetTargetPath("")
	cmd.Stdin = nil
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	if !step.Quiet {
		outs := stdout.String()
		if len(outs) > 0 {
			exec.PrintInfo(outs)
		}
	}
	errs := stderr.String()
	if len(errs) > 0 {
		exec.PrintError(errs)
	}

	return err
}

func (step ShellStep) getShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	return shell
}
