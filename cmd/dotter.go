package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/jayclassless/dotter"
)

var (
	version = "dev"

	app = kingpin.New(
		"dotter",
		"Installs a collection of dotfiles into a directory.",
	)

	sourcePath = app.Arg(
		"source",
		"Path to the dotfile collection to install.",
	).String()

	targetPath = app.Arg(
		"target",
		"Path to install the dotfiles to.",
	).String()

	quiet = app.Flag(
		"quiet",
		"Surpress all output from dotter.",
	).Short('q').Bool()

	continueOnError = app.Flag(
		"continue-on-error",
		"Continue execution even if a step fails.",
	).Short('c').Bool()
)

func cleanPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[1:])
	}

	path, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", err
	}

	return path, nil
}

func determineSource(path string) (string, string, error) {
	var cfg string
	var err error

	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			return "", "", err
		}
	}

	path, err = cleanPath(path)
	if err != nil {
		return "", "", err
	}

	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", "", err
	}

	if fileInfo.IsDir() {
		cfg = filepath.Join(path, "dotter.yaml")
	} else if fileInfo.Mode().IsRegular() {
		cfg = path
		path = filepath.Dir(path)
	} else {
		return "", "", fmt.Errorf("%s is not a valid source", path)
	}

	return path, cfg, nil
}

func determineTarget(path string) (string, error) {
	var err error

	if path == "" {
		path, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}

	path, err = cleanPath(path)
	if err != nil {
		return "", err
	}

	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", err
	}
	if !fileInfo.IsDir() {
		return "", fmt.Errorf("%s is not a directory", path)
	}

	return path, nil
}

func failIfError(err error, message string) {
	if err != nil {
		app.FatalIfError(err, message)
	}
}

func main() {
	app.Version(version)
	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	sourcePath, configPath, err := determineSource(*sourcePath)
	failIfError(err, "Could not determine source path")
	targetPath, err := determineTarget(*targetPath)
	failIfError(err, "Could not determine target path")

	config, err := dotter.NewConfigurationFromFile(configPath)
	failIfError(err, "Could not read configuration file")
	config.Options.Quiet = *quiet
	config.Options.StopOnError = !*continueOnError

	exec := dotter.NewExecutor(sourcePath, targetPath, config)
	err = exec.Execute()
	if err != nil {
		os.Exit(1)
	}
}
