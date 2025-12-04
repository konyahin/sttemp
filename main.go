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

type TemplateFile struct {
	Name    string
	Content []byte
}

func findTemplates(baseDir string) ([]TemplateFile, error) {
	var templates []TemplateFile
	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && (strings.Contains(d.Name(), "/.") || strings.HasPrefix(d.Name(), ".")) {
			return fs.SkipDir
		}

		if !d.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			relPath, err := filepath.Rel(baseDir, path)
			if err != nil {
				return err
			}

			templates = append(templates, TemplateFile{
				relPath,
				content,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return templates, nil
}

func main() {
	path := flag.String("C", "", "template's directory")
	flag.Parse()

	baseDir, err := getBaseDir(*path)
	if err != nil {
		log.Fatal(err)
	}

	templates, err := findTemplates(baseDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range templates {
		template := NewTemplate(t.Name, t.Content)
		fmt.Println(template)
	}
}
