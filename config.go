package dotter

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v3"

	"github.com/jayclassless/dotter/step"
)

type StepDefaultOptions struct {
	Link      step.LinkOptions
	Directory step.DirectoryOptions
	Shell     step.ShellOptions
	Clean     step.CleanOptions
}

func NewStepDefaultOptions() StepDefaultOptions {
	opt := StepDefaultOptions{}
	opt.Link = step.NewLinkOptions()
	opt.Directory = step.NewDirectoryOptions()
	opt.Shell = step.NewShellOptions()
	opt.Clean = step.NewCleanOptions()
	return opt
}

type Options struct {
	BackupForced string
	StopOnError  bool
	Quiet        bool
	Defaults     StepDefaultOptions
}

func NewOptions() Options {
	options := Options{}
	options.StopOnError = true
	options.Quiet = false
	options.Defaults = NewStepDefaultOptions()
	return options
}

type Configuration struct {
	SourcePath string
	Options    Options
	Steps      []step.Step
}

func NewConfiguration() Configuration {
	cfg := Configuration{}
	cfg.Options = NewOptions()
	cfg.Steps = make([]step.Step, 0)
	return cfg
}

func NewConfigurationFromFile(configPath string) (Configuration, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return NewConfiguration(), err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return NewConfiguration(), err
	}

	cfg, err := NewConfigurationFromYaml(content)
	cfg.SourcePath = configPath

	return cfg, err
}

type yamlConfig struct {
	Options Options
	Steps   []yaml.Node
}

func NewConfigurationFromYaml(content []byte) (Configuration, error) {
	cfg := NewConfiguration()

	tmpCfg := yamlConfig{}
	tmpCfg.Options = NewOptions()
	err := yaml.Unmarshal(content, &tmpCfg)
	if err != nil {
		return cfg, err
	}
	cfg.Options = tmpCfg.Options

	for _, node := range tmpCfg.Steps {
		if node.Kind == yaml.MappingNode {
			steps, err := parseStepsFromNode(node, cfg.Options.Defaults)
			if err != nil {
				return cfg, err
			}
			cfg.Steps = append(cfg.Steps, steps...)
		} else {
			return cfg, fmt.Errorf("Unexpected %s value at line %d", node.Tag, node.Line)
		}
	}

	return cfg, nil
}

func parseStepsFromNode(node yaml.Node, defaults StepDefaultOptions) ([]step.Step, error) {
	stepName := node.Content[0].Value

	if stepName == "link" {
		return parseLinkBlock(node.Content[1], defaults.Link)
	} else if stepName == "directory" {
		return parseDirectoryBlock(node.Content[1], defaults.Directory)
	} else if stepName == "shell" {
		return parseShellBlock(node.Content[1], defaults.Shell)
	} else if stepName == "clean" {
		return parseCleanBlock(node.Content[1], defaults.Clean)
	} else if stepName == "include_steps" {
		// TODO
		return nil, nil
	}

	return nil, fmt.Errorf("Unexpected step type \"%s\"", stepName)
}

func parseLinkBlock(node *yaml.Node, defaults step.LinkOptions) ([]step.Step, error) {
	steps := make([]step.Step, 0)

	if node.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("Link definitions not in a mapping at line %d", node.Line)
	}
	nodes := node.Content

	for i := 0; i < len(nodes); i += 2 {
		link := step.NewLinkStepWithDefaults(defaults)
		link.Target = nodes[i].Value

		details := nodes[i+1]
		if details.Tag == "!!str" {
			link.Source = details.Value

		} else if details.Kind == yaml.MappingNode {
			err := details.Decode(&link)
			if err != nil {
				return nil, err
			}

		} else {
			return nil, fmt.Errorf("Unexpected link definition type %s at line %d", details.Tag, details.Line)
		}

		steps = append(steps, link)
	}

	return steps, nil
}

func parseDirectoryBlock(node *yaml.Node, defaults step.DirectoryOptions) ([]step.Step, error) {
	steps := make([]step.Step, 0)

	if node.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("Directory definitions not in a sequence at line %d", node.Line)
	}

	for _, details := range node.Content {
		dir := step.NewDirectoryStepWithDefaults(defaults)

		if details.Tag == "!!str" {
			dir.Target = details.Value

		} else if details.Kind == yaml.MappingNode {
			err := details.Decode(&dir)
			if err != nil {
				return nil, err
			}

		} else {
			return nil, fmt.Errorf("Unexpected directory definition type %s at line %d", details.Tag, details.Line)
		}

		steps = append(steps, dir)
	}

	return steps, nil
}

func parseShellBlock(node *yaml.Node, defaults step.ShellOptions) ([]step.Step, error) {
	steps := make([]step.Step, 0)

	if node.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("Shell definitions not in a sequence at line %d", node.Line)
	}

	for _, details := range node.Content {
		shell := step.NewShellStepWithDefaults(defaults)

		if details.Tag == "!!str" {
			shell.Command = details.Value

		} else if details.Kind == yaml.MappingNode {
			err := details.Decode(&shell)
			if err != nil {
				return nil, err
			}

		} else {
			return nil, fmt.Errorf("Unexpected shell definition type %s at line %d", details.Tag, details.Line)
		}

		steps = append(steps, shell)
	}

	return steps, nil
}

func parseCleanBlock(node *yaml.Node, defaults step.CleanOptions) ([]step.Step, error) {
	steps := make([]step.Step, 0)

	if node.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("Clean definitions not in a sequence at line %d", node.Line)
	}

	for _, details := range node.Content {
		clean := step.NewCleanStepWithDefaults(defaults)

		if details.Tag == "!!str" {
			clean.Target = details.Value

		} else if details.Kind == yaml.MappingNode {
			err := details.Decode(&clean)
			if err != nil {
				return nil, err
			}

		} else {
			return nil, fmt.Errorf("Unexpected clean definition type %s at line %d", details.Tag, details.Line)
		}

		steps = append(steps, clean)
	}

	return steps, nil
}
