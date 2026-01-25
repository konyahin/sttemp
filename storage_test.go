package main

import (
	"errors"
	"io/fs"
	"testing"
)

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
			"when path is relative it shoul be turn into absolute path",
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
			"if path is empty UserHomeDir should not affect storage",
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

			if err == nil && storage.path != tt.wantPath {
				t.Fatalf("storage directory should be \"%v\", but we got \"%v\"", tt.wantPath, storage.path)
			}
		})
	}
}
