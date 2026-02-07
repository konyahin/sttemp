package main

import (
	"os/exec"
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
			},
			wantErr: "both -d and -o flags were set, but only one of them can be used at the same time",
		},
		{
			name: "edit mode without template name",
			clistate: CliState{
				outputFileName: "",
				templateNames:  []string{},
				storage:        &Storage{templates: map[string]TemplateFile{}},
				editMode:       true,
			},
			wantErr: "edit mode was set, but no template name was provided",
		},
		{
			name: "edit mode with too many template names",
			clistate: CliState{
				outputFileName: "",
				templateNames:  []string{"template1", "template2"},
				storage:        &Storage{templates: map[string]TemplateFile{}},
				editMode:       true,
			},
			wantErr: "edit mode was set, but too many template names were provided",
		},
		{
			name: "non-existent template",
			clistate: CliState{
				outputFileName: "",
				templateNames:  []string{"nonexistent"},
				storage:        &Storage{templates: map[string]TemplateFile{}},
			},
			wantErr: "template nonexistent not found",
		},
		{
			name: "template with no default name but -d flag set",
			clistate: CliState{
				outputFileName: "",
				defaultName:    true,
				templateNames:  []string{"template-without-default"},
				storage: &Storage{templates: map[string]TemplateFile{
					"template-without-default": {
						Name:        "template-without-default",
						DefaultName: "",
						Path:        "/path/to/template",
					},
				}},
			},
			wantErr: "template template-without-default has no default name, but -d flag was set",
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

func TestEditMode(t *testing.T) {
	t.Run("edit mode should use EDITOR environment variable", func(t *testing.T) {
		var commandName string
		var commandArgs []string
		ioh := &IOHandler{
			LookupEnv: func(key string) (string, bool) {
				if key == "EDITOR" {
					return "emacs", true
				}
				return "", false
			},
			Command: func(name string, arg ...string) *exec.Cmd {
				commandName = name
				commandArgs = arg
				return &exec.Cmd{
					Path: "/usr/bin/true",
				}
			},
		}
		cliState := CliState{
			templateNames: []string{"for-edit"},
			storage: &Storage{templates: map[string]TemplateFile{
				"for-edit": {
					Name:        "for-edit",
					DefaultName: "",
					Path:        "/path/to/template/for-edit",
				},
			}},
			ioh:      ioh,
			editMode: true,
		}

		err := cliState.Run()
		if err != nil {
			t.Fatalf("edit mode should not return errors, but we got %v", err)
		}

		if commandName != "emacs" {
			t.Fatalf("edit mode should use EDITOR var, but we got %v", commandName)
		}

		if len(commandArgs) != 1 || commandArgs[0] != "/path/to/template/for-edit" {
			t.Fatalf("edit mode should use template path as editor argument, but we got %v", commandArgs)
		}
	})

	t.Run("edit mode should use vi, if EDITOR environment variable is missing", func(t *testing.T) {
		var commandName string
		var commandArgs []string
		ioh := &IOHandler{
			LookupEnv: func(key string) (string, bool) {
				return "", false
			},
			Command: func(name string, arg ...string) *exec.Cmd {
				commandName = name
				commandArgs = arg
				return &exec.Cmd{
					Path: "/usr/bin/true",
				}
			},
		}
		cliState := CliState{
			templateNames: []string{"for-edit"},
			storage: &Storage{templates: map[string]TemplateFile{
				"for-edit": {
					Name:        "for-edit",
					DefaultName: "",
					Path:        "/path/to/template/for-edit",
				},
			}},
			ioh:      ioh,
			editMode: true,
		}

		err := cliState.Run()
		if err != nil {
			t.Fatalf("edit mode should not return errors, but we got %v", err)
		}

		if commandName != "vi" {
			t.Fatalf("edit mode should use vi, but we got %v", commandName)
		}

		if len(commandArgs) != 1 || commandArgs[0] != "/path/to/template/for-edit" {
			t.Fatalf("edit mode should use template path as editor argument, but we got %v", commandArgs)
		}
	})

	t.Run("edit mode should return errors, if something goes wrong", func(t *testing.T) {
		ioh := &IOHandler{
			LookupEnv: func(key string) (string, bool) {
				return "", false
			},
			Command: func(name string, arg ...string) *exec.Cmd {
				return &exec.Cmd{
					Path: "/usr/bin/false",
				}
			},
		}
		cliState := CliState{
			templateNames: []string{"for-edit"},
			storage: &Storage{templates: map[string]TemplateFile{
				"for-edit": {
					Name:        "for-edit",
					DefaultName: "",
					Path:        "/path/to/template/for-edit",
				},
			}},
			ioh:      ioh,
			editMode: true,
		}

		err := cliState.Run()
		if err == nil || err.Error() != "exit status 1" {
			t.Fatalf("edit mode should return error , but we got %v", err)
		}
	})
}
