package main

import (
	"errors"
	"io/fs"
	"os"
	"reflect"
	"testing"
)

type Dir struct {
	// path to file/dir in file system
	path  string
	// name of the file/dir
	name  string
	// is it dir?
	isDir bool
	// error from the system when we try to access file
	err   error
}

func (d Dir) Name() string {
	return d.name
}

func (d Dir) IsDir() bool {
	return d.isDir
}

// should be unused in code
func (d Dir) Type() fs.FileMode {
	return 0
}

// should be unused in code
func (d Dir) Info() (fs.FileInfo, error) {
	return nil, nil
}

func TestStorageTemplateDir(t *testing.T) {
	testCases := []struct {
		name           string
		path           string
		userHomeDirErr error
		wantError      error
		wantPath       string
	}{
		{
			"when path is empty, use user home dir",
			"",
			nil,
			nil,
			"HOME_DIR/.local/share/sttemp",
		},
		{
			"when path is set, use it",
			"/usr/local/sttemp",
			nil,
			nil,
			"/usr/local/sttemp",
		},
		{
			"when path is relative it should be turned into absolute path",
			"/usr/local/../sttemp",
			nil,
			nil,
			"/usr/sttemp",
		},
		{
			"if path is empty and UserHomeDir return error, we should got it",
			"",
			fs.ErrPermission,
			fs.ErrPermission,
			"",
		},
		{
			"if path is not empty UserHomeDir should not affect storage",
			"/usr/local/sttemp",
			fs.ErrPermission,
			nil,
			"/usr/local/sttemp",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ioh := &IOHandler{
				UserHomeDir: func() (string, error) {
					return "HOME_DIR", tt.userHomeDirErr
				},
				WalkDir: func(root string, fn fs.WalkDirFunc) error {
					return nil
				},
			}

			storage, err := NewStorage(ioh, tt.path)
			if !errors.Is(err, tt.wantError) {
				t.Fatalf("err should be \"%v\", but we got \"%v\"", tt.wantError, err)
			}

			if tt.wantError != nil && storage != nil {
				t.Error("storage should be nil when error is returned")
			}

			if tt.wantError == nil && storage.path != tt.wantPath {
				t.Fatalf("storage directory should be \"%v\", but we got \"%v\"", tt.wantPath, storage.path)
			}
		})
	}
}

func TestStorageTemplates(t *testing.T) {
	testCases := []struct {
		name   string
		walk   []Dir
		expect map[string]TemplateFile
		err    error
	}{
		{
			"empty directory",
			[]Dir{},
			map[string]TemplateFile{},
			nil,
		},
		{
			"happy path",
			[]Dir{
				{
					"/templates/first",
					"first",
					false,
					nil,
				},
				{
					"/templates/LICENSE",
					"LICENSE",
					true,
					nil,
				},
				{
					"/templates/LICENSE/mit",
					"mit",
					false,
					nil,
				},
				{
					"/templates/LICENSE/gpl",
					"gpl",
					false,
					nil,
				},
			},
			map[string]TemplateFile{
				"first": {
					"first",
					"",
					"/templates/first",
				},
				"mit": {
					"mit",
					"LICENSE",
					"/templates/LICENSE/mit",
				},
				"gpl": {
					"gpl",
					"LICENSE",
					"/templates/LICENSE/gpl",
				},
			},
			nil,
		},
		{
			"skip hidden dir",
			[]Dir{
				{
					"/templates/.config",
					".config",
					true,
					nil,
				},
			},
			nil,
			fs.SkipDir,
		},
		{
			"return fs errors as is (except permissions errors)",
			[]Dir{
				{
					"/templates/first",
					"first",
					false,
					fs.ErrInvalid,
				},
			},
			nil,
			fs.ErrInvalid,
		},
		{
			"ignore permissions errors",
			[]Dir{
				{
					"/templates/first",
					"first",
					false,
					nil,
				},
				{
					"/templates/second",
					"second",
					false,
					os.ErrPermission,
				},
			},
			map[string]TemplateFile{
				"first": {
					"first",
					"",
					"/templates/first",
				},
			},
			nil,
		},
		{
			"two templates with the same name",
			[]Dir{
				{
					"/templates/first",
					"first",
					false,
					nil,
				},
				{
					"/templates/subdir/first",
					"first",
					false,
					nil,
				},
			},
			nil,
			ErrDuplicateTemplate,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ioh := &IOHandler{
				WalkDir: func(root string, fn fs.WalkDirFunc) error {
					for _, walk := range tt.walk {
						err := fn(walk.path, walk, walk.err)
						if err != nil {
							return err
						}
					}
					return nil
				},
			}

			storage, err := NewStorage(ioh, "/templates")

			if !errors.Is(err, tt.err) {
				t.Fatalf("err should be \n%v\nbut we got\n%v\n", tt.err, err)
			}

			if tt.err == nil && !reflect.DeepEqual(tt.expect, storage.templates) {
				t.Fatalf("Storage should contain\n%#v\nbut we got %#v\n", tt.expect, storage.templates)
			}
		})
	}
}
