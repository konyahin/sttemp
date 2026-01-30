package main

import (
	"testing"
)

func TestCliValidation(t *testing.T) {
	testCases := []struct {
		name     string
		clistate CliState
		wantErr  string
	}{
		{
			name: "conflicting -d and -o flags",
			clistate: CliState{
				outputFileName: "output.txt",
				defaultName:    true,
				templateNames:  []string{},
				storage:        &Storage{templates: map[string]TemplateFile{}},
				noInput:        false,
				ioh:            nil,
				editMode:       false,
				listTemplates:  false,
			},
			wantErr: "both -d and -o flags were set, but only one of them can be used at the same time",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.clistate.Run()

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("expected no error, but got: %v", err)
				}
			} else {
				if err == nil {
					t.Fatalf("expected error %q, but got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("expected error:\n%q\nbut got:\n%q", tt.wantErr, err.Error())
				}
			}
		})
	}
}
