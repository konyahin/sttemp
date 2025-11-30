package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const defaultTemplateDir = ".local/share/sttemp"

type State struct {
	TemplateDir string
}

func (s *State) DefaultTemplateDir() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	s.TemplateDir = filepath.Join(home, defaultTemplateDir)
}

func (s *State) SetTemplateDir(templatesPath string) {
	templatesPath, err := filepath.Abs(templatesPath)
	if err != nil {
		log.Fatal(err)
	}

	s.TemplateDir = templatesPath
}

func (s State) Templates() []*Template {
	var templates []*Template

	err := filepath.WalkDir(s.TemplateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && (strings.Contains(d.Name(), "/.") || strings.HasPrefix(d.Name(), ".")) {
			return fs.SkipDir
		}

		if !d.IsDir() {
			templates = append(templates, NewTemplate(s, path))
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return templates
}
