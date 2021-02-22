package step

import "fmt"

type CleanOptions struct {
	Force     bool
	Recursive bool
}

func NewCleanOptions() CleanOptions {
	opt := CleanOptions{}
	opt.Force = false
	opt.Recursive = false
	return opt
}

type CleanStep struct {
	CleanOptions `yaml:",inline"`
	Target       string `yaml:"path"`
}

func NewCleanStep() CleanStep {
	return NewCleanStepWithDefaults(NewCleanOptions())
}

func NewCleanStepWithDefaults(defaults CleanOptions) CleanStep {
	step := CleanStep{}
	step.CleanOptions = defaults
	return step
}

func (step CleanStep) GetActivityLabel() string {
	return "Cleaning"
}

func (step CleanStep) GetActivityDetails() string {
	return step.Target
}

func (step CleanStep) Execute(exec StepExecutor) error {
	fmt.Printf("Cleaning %s\n", step.Target)
	return nil
}
