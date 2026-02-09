package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CommandRunner interface {
	Run(ioh *IOHandler, name string, args ...string) error
}

type RealCommandRunner struct{}

func (r *RealCommandRunner) Run(ioh *IOHandler, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = ioh.Stdin
	cmd.Stdout = ioh.Stdout
	cmd.Stderr = ioh.Stderr
	return cmd.Run()
}

type IOHandler struct {
	Stdin         io.Reader
	Stdout        io.Writer
	Stderr        io.Writer
	LookupEnv     func(key string) (string, bool)
	ReadFile      func(name string) ([]byte, error)
	UserHomeDir   func() (string, error)
	WalkDir       func(root string, fn fs.WalkDirFunc) error
	Create        func(name string) (OutputFile, error)
	CommandRunner CommandRunner
}

func DefaultIOHandler() *IOHandler {
	return &IOHandler{
		Stdin:       os.Stdin,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		LookupEnv:   os.LookupEnv,
		ReadFile:    os.ReadFile,
		UserHomeDir: os.UserHomeDir,
		WalkDir:     filepath.WalkDir,
		Create: func(name string) (OutputFile, error) {
			file, err := os.Create(name)
			return OutputFile(file), err
		},
		CommandRunner: &RealCommandRunner{},
	}
}

func (ioh *IOHandler) askForValue(variable string) (string, error) {
	reader := bufio.NewReader(ioh.Stdin)
	fmt.Fprintf(ioh.Stderr, "Enter value for %s: ", variable)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimRight(input, "\n"), nil
}

func (ioh *IOHandler) getVariableValues(template *Template, noInput bool) (map[string]string, error) {
	values := make(map[string]string, len(template.Variables))
	for _, variable := range template.Variables {
		envValue, ok := ioh.LookupEnv(variable)
		if !ok && noInput {
			return nil, fmt.Errorf("variable %s is not set and --no-input is enabled; set %s in environment", variable, variable)
		}
		if !ok {
			var err error
			envValue, err = ioh.askForValue(variable)
			if err != nil {
				return nil, err
			}
		}
		values[variable] = envValue
	}
	return values, nil
}

func (ioh *IOHandler) executeCommand(command string, arg string) error {
	// command can be a single command or a command with arguments
	// if it has arguments, split it and add to args
	var args []string
	if strings.Contains(command, " ") {
		args = strings.Split(command, " ")
		args = append(args, arg)
	} else {
		args = []string{command, arg}
	}

	return ioh.CommandRunner.Run(ioh, args[0], args[1:]...)
}
