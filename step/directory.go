package step

import (
	"fmt"
	"os"
	"path/filepath"
)

// DirectoryOptions contains non-path options for Directory steps
type DirectoryOptions struct {
	CreateParents bool `yaml:"create_parents"`
	Mode          uint
	Force         bool
}

// NewDirectoryOptions creates a new instance of a DirectoryOptions struct
func NewDirectoryOptions() DirectoryOptions {
	opt := DirectoryOptions{}
	opt.CreateParents = true
	opt.Mode = 0o755
	opt.Force = false
	return opt
}

// DirectoryStep contains the specification for Directory steps
type DirectoryStep struct {
	DirectoryOptions `yaml:",inline"`
	Target           string `yaml:"path"`
}

// NewDirectoryStep creates a new instance of a DirectoryStep struct using default options
func NewDirectoryStep() DirectoryStep {
	return NewDirectoryStepWithDefaults(NewDirectoryOptions())
}

// NewDirectoryStepWithDefaults creates a new instsance of a DirectoryStep struct using the specified options
func NewDirectoryStepWithDefaults(defaults DirectoryOptions) DirectoryStep {
	step := DirectoryStep{}
	step.DirectoryOptions = defaults
	return step
}

// GetActivityLabel returns a short description of what a DirectoryStep does
func (step DirectoryStep) GetActivityLabel() string {
	return "Directory"
}

// GetActivityDetails returns description specific to this particular instance of the DirectoryStep
func (step DirectoryStep) GetActivityDetails() string {
	return step.Target
}

// Execute creates the specified directory
func (step DirectoryStep) Execute(exec StepExecutor) error {
	targetPath := exec.GetTargetPath(step.Target)

	fileInfo, err := os.Stat(targetPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if !fileInfo.IsDir() {
			if step.Force {
				err = exec.ForceRemove(targetPath)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("Non-directory %s already exists", targetPath)
			}
		}
	}

	desiredMode := os.FileMode(step.Mode)

	fileInfo, err = os.Stat(targetPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		parentPath := filepath.Dir(targetPath)
		_, err := os.Stat(parentPath)
		if os.IsNotExist(err) {
			if step.CreateParents {
				os.MkdirAll(targetPath, desiredMode)
			} else {
				return fmt.Errorf(
					"Cannot create %s as parent directory %s does not exist",
					step.Target,
					parentPath,
				)
			}
		} else if err == nil {
			os.Mkdir(targetPath, desiredMode)
		} else {
			return err
		}

		fileInfo, err = os.Stat(targetPath)
	}

	if fileInfo.Mode().Perm() != desiredMode.Perm() {
		os.Chmod(targetPath, desiredMode)
	}

	return nil
}
