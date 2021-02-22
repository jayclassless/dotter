package step

import (
	"fmt"
	"os"
	"path/filepath"
)

// LinkOptions contains non-path options for Link steps
type LinkOptions struct {
	CreateParents bool `yaml:"create_parents"`
	Relative      bool
	Force         bool
	Relink        bool
}

// NewLinkOptions creates a new instance of a LinkOptions struct
func NewLinkOptions() LinkOptions {
	opt := LinkOptions{}
	opt.CreateParents = true
	opt.Relative = true
	opt.Force = false
	opt.Relink = true
	return opt
}

// LinkStep contains the specification for Link steps
type LinkStep struct {
	LinkOptions `yaml:",inline"`
	Target      string
	Source      string
}

// NewLinkStep creates a new instance of a LinkStep struct using default options
func NewLinkStep() LinkStep {
	return NewLinkStepWithDefaults(NewLinkOptions())
}

// NewLinkStepWithDefaults creates a new instsance of a LinkStep struct using the specified options
func NewLinkStepWithDefaults(defaults LinkOptions) LinkStep {
	step := LinkStep{}
	step.LinkOptions = defaults
	return step
}

// GetActivityLabel returns a short description of what a LinkStep does
func (step LinkStep) GetActivityLabel() string {
	return "Linking"
}

// GetActivityDetails returns description specific to this particular instance of the LinkStep
func (step LinkStep) GetActivityDetails() string {
	return step.Target
}

// Execute creates the specified symlink
func (step LinkStep) Execute(exec StepExecutor) error {
	var err error

	targetPath := exec.GetTargetPath(step.Target)
	sourcePath := exec.GetSourcePath(step.Source)
	if step.Relative {
		sourcePath, err = filepath.Rel(filepath.Dir(targetPath), sourcePath)
		if err != nil {
			return err
		}
	}
	parentPath := filepath.Dir(targetPath)

	fileInfo, err := os.Lstat(targetPath)
	if err == nil {
		if IsSymLink(fileInfo) {
			current, err := os.Readlink(targetPath)
			if err != nil {
				return nil
			}
			if current == sourcePath {
				// Link exists and is pointing to the right thing
				return nil

			} else if step.Relink {
				// Link exists, but is wrong, and we want to fix it
				err = os.Remove(targetPath)
				if err != nil {
					return err
				}
				return os.Symlink(sourcePath, targetPath)

			}

			// Link exists, but is wrong
			return fmt.Errorf(
				"Cannot create %s as a symlink because one already exists",
				targetPath,
			)
		}

		if step.Force {
			// Something other than a link exists, and we want to replace it
			err = exec.ForceRemove(targetPath)
			if err != nil {
				return err
			}
			return os.Symlink(sourcePath, targetPath)
		}

		// Something other than a link exists
		return fmt.Errorf("Non-link %s already exists", targetPath)

	} else if os.IsNotExist(err) {
		// Nothing exists, make the link
		_, err := os.Stat(parentPath)
		if os.IsNotExist(err) {
			if step.CreateParents {
				// Parent dir doesn't exist, make it first
				os.MkdirAll(parentPath, os.FileMode(0o777))
			} else {
				// Parent dir doesn't exist
				return fmt.Errorf(
					"Cannot create %s as parent directory %s does not exist",
					step.Target,
					parentPath,
				)
			}
		} else if err != nil {
			return err
		}

		return os.Symlink(sourcePath, targetPath)
	}

	return err
}
