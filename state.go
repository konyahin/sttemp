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

func NewStateWithDefaultDir() *State {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(home, defaultTemplateDir)
	return NewState(path)
}

func NewState(templatesPath string) *State {
	templatesPath, err := filepath.Abs(templatesPath)
	if err != nil {
		log.Fatal(err)
	}

	state := new(State)
	state.TemplateDir = templatesPath
	return state
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
