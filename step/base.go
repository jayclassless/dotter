package step

// StepExecutor defines the interface necessary to run Step.Execute()
type StepExecutor interface {
	GetTargetPath(path string) string
	GetSourcePath(path string) string
	ForceRemove(path string) error
	PrintInfo(message string)
	PrintError(message string)
}

// Step defines the interface necessary for an installation step
type Step interface {
	GetActivityLabel() string
	GetActivityDetails() string
	Execute(StepExecutor) error
}
