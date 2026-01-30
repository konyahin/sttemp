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
	path string
	// name of the file/dir
	name string
	// is it dir?
	isDir bool
	// error from the system when we try to access file
	err error
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
			name:           "when path is empty, use user home dir",
			path:           "",
			userHomeDirErr: nil,
			wantError:      nil,
			wantPath:       "HOME_DIR/.local/share/sttemp",
		},
		{
			name:           "when path is set, use it",
			path:           "/usr/local/sttemp",
			userHomeDirErr: nil,
			wantError:      nil,
			wantPath:       "/usr/local/sttemp",
		},
		{
			name:           "when path is relative it should be turned into absolute path",
			path:           "/usr/local/../sttemp",
			userHomeDirErr: nil,
			wantError:      nil,
			wantPath:       "/usr/sttemp",
		},
		{
			name:           "if path is empty and UserHomeDir return error, we should got it",
			path:           "",
			userHomeDirErr: fs.ErrPermission,
			wantError:      fs.ErrPermission,
			wantPath:       "",
		},
		{
			name:           "if path is not empty UserHomeDir should not affect storage",
			path:           "/usr/local/sttemp",
			userHomeDirErr: fs.ErrPermission,
			wantError:      nil,
			wantPath:       "/usr/local/sttemp",
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
			name:   "empty directory",
			walk:   []Dir{},
			expect: map[string]TemplateFile{},
			err:    nil,
		},
		{
			name: "happy path",
			walk: []Dir{
				{
					path:  "/templates/first",
					name:  "first",
					isDir: false,
					err:   nil,
				},
				{
					path:  "/templates/LICENSE",
					name:  "LICENSE",
					isDir: true,
					err:   nil,
				},
				{
					path:  "/templates/LICENSE/mit",
					name:  "mit",
					isDir: false,
					err:   nil,
				},
				{
					path:  "/templates/LICENSE/gpl",
					name:  "gpl",
					isDir: false,
					err:   nil,
				},
			},
			expect: map[string]TemplateFile{
				"first": {
					Name:        "first",
					DefaultName: "",
					Path:        "/templates/first",
				},
				"mit": {
					Name:        "mit",
					DefaultName: "LICENSE",
					Path:        "/templates/LICENSE/mit",
				},
				"gpl": {
					Name:        "gpl",
					DefaultName: "LICENSE",
					Path:        "/templates/LICENSE/gpl",
				},
			},
			err: nil,
		},
		{
			name: "skip hidden dir",
			walk: []Dir{
				{
					path:  "/templates/.config",
					name:  ".config",
					isDir: true,
					err:   nil,
				},
			},
			expect: nil,
			err:    fs.SkipDir,
		},
		{
			name: "return fs errors as is (except permissions errors)",
			walk: []Dir{
				{
					path:  "/templates/first",
					name:  "first",
					isDir: false,
					err:   fs.ErrInvalid,
				},
			},
			expect: nil,
			err:    fs.ErrInvalid,
		},
		{
			name: "ignore permissions errors",
			walk: []Dir{
				{
					path:  "/templates/first",
					name:  "first",
					isDir: false,
					err:   nil,
				},
				{
					path:  "/templates/second",
					name:  "second",
					isDir: false,
					err:   os.ErrPermission,
				},
			},
			expect: map[string]TemplateFile{
				"first": {
					Name:        "first",
					DefaultName: "",
					Path:        "/templates/first",
				},
			},
			err: nil,
		},
		{
			name: "two templates with the same name",
			walk: []Dir{
				{
					path:  "/templates/first",
					name:  "first",
					isDir: false,
					err:   nil,
				},
				{
					path:  "/templates/subdir/first",
					name:  "first",
					isDir: false,
					err:   nil,
				},
			},
			expect: nil,
			err:    ErrDuplicateTemplate,
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
