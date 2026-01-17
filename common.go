package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getVariableValues(template *Template, noInput bool) (map[string]string, error) {
	values := make(map[string]string, len(template.Variables))
	for _, variable := range template.Variables {
		envValue, ok := os.LookupEnv(variable)
		if !ok && noInput {
			return nil, fmt.Errorf("variable %s is not set and --no-input is enabled; set %s in environment", variable, variable)
		}
		if !ok {
			var err error
			envValue, err = askForValue(variable)
			if err != nil {
				return nil, err
			}
		}
		values[variable] = envValue
	}
	return values, nil
}

func askForValue(variable string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter value for %s: ", variable)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimRight(input, "\n"), nil
}
