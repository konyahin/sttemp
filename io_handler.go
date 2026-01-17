package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type IOHandler struct {
	io.Reader
	io.Writer
	LookupEnv   func(key string) (string, bool)
	ReadFile    func(name string) ([]byte, error)
	UserHomeDir func() (string, error)
	WalkDir     func(root string, fn fs.WalkDirFunc) error
}

func DefaultIOHandler() *IOHandler {
	return &IOHandler{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		LookupEnv:   os.LookupEnv,
		ReadFile:    os.ReadFile,
		UserHomeDir: os.UserHomeDir,
		WalkDir:     filepath.WalkDir,
	}
}

func (ioh *IOHandler) askForValue(variable string) (string, error) {
	reader := bufio.NewReader(ioh)
	fmt.Fprintf(ioh, "Enter value for %s: ", variable)

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

func (ioh *IOHandler) create(name string) (OutputFile, error) {
	return os.Create(name)
}
