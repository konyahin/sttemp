package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Storage represent a directory with all templates
type Storage struct {
	path      string
	templates map[string]TemplateFile
}

func NewStorage(ioh *IOHandler, path string) (*Storage, error) {
	path, err := getStoragePath(ioh, path)
	if err != nil {
		return nil, err
	}

	templateFiles, err := findTemplateFiles(ioh, path)
	if err != nil {
		return nil, err
	}

	storage := &Storage{
		path:      path,
		templates: templateFiles,
	}
	return storage, nil
}

func GetDefaultTemplateDir() string {
	return ".local/share/sttemp"
}

func getStoragePath(ioh *IOHandler, path string) (string, error) {
	if path == "" {
		home, err := ioh.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, GetDefaultTemplateDir()), nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

func findTemplateFiles(ioh *IOHandler, path string) (map[string]TemplateFile, error) {
	templateFiles := make(map[string]TemplateFile)
	err := ioh.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return nil
			}
			return err
		}

		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return fs.SkipDir
		}

		if !d.IsDir() {
			templateFile, err := NewTemplateFile(filePath, path)
			if err != nil {
				return err
			}
			if old, ok := templateFiles[templateFile.Name]; ok {
				return fmt.Errorf("two templates with same file name:\n\t%v\n\t%v", old.Path, templateFile.Path)
			}
			templateFiles[templateFile.Name] = *templateFile
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return templateFiles, nil
}
