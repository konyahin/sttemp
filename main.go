package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const defaultTemplateDir = ".local/share/sttemp"

func getDefaultBaseDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, defaultTemplateDir), nil
}

func getBaseDir(baseDir string) (string, error) {
	if baseDir == "" {
		return getDefaultBaseDir()
	}

	templatesPath, err := filepath.Abs(baseDir)
	if err != nil {
		return "", nil
	}

	return templatesPath, nil
}

func findTemplateFiles(baseDir string) ([]*TemplateFile, error) {
	var templateFiles []*TemplateFile
	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && (strings.Contains(d.Name(), "/.") || strings.HasPrefix(d.Name(), ".")) {
			return fs.SkipDir
		}

		if !d.IsDir() {
			templateFiles = append(templateFiles, NewTemplateFile(path, baseDir))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return templateFiles, nil
}

func main() {
	path := flag.String("C", "", "template's directory")
	flag.Parse()

	baseDir, err := getBaseDir(*path)
	if err != nil {
		log.Fatal(err)
	}

	templateFiles, err := findTemplateFiles(baseDir)
	if err != nil {
		log.Fatal(err)
	}

	names := flag.Args()
	if len(names) > 0 {
		for _, name := range names {
			for _, templateFile := range templateFiles {
				if name != templateFile.Name {
					continue
				}

				content, err := os.ReadFile(templateFile.Path)
				if err != nil {
					log.Fatal(err)
				}

				template := NewTemplate(templateFile, content)
				values := make(map[string]string)
				for _, variable := range template.Variables {
					values[variable] = os.Getenv(variable)
				}
				fmt.Println(template.fillTemplate(values))
			}
		}
		return
	}

	// if no args, print all template names
	for _, templateFile := range templateFiles {
		fmt.Println(templateFile)
	}
}
