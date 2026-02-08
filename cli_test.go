package main

import (
	"bytes"
	"errors"
	"maps"
	"slices"
	"testing"
)

type MockCommandRunner struct {
	ShouldFail   bool
	CapturedName string
	CapturedArgs []string
}

func (m *MockCommandRunner) Run(ioh *IOHandler, name string, args ...string) error {
	m.CapturedName = name
	m.CapturedArgs = args

	if m.ShouldFail {
		return errors.New("simulated error")
	}

	return nil
}

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
	testCases := []struct {
		name        string
		editorVar   string
		shouldFail  bool
		expectedCmd string
		expectedArg string
		expectError bool
	}{
		{
			name:        "edit mode should use EDITOR environment variable",
			editorVar:   "emacs",
			expectedCmd: "emacs",
			expectedArg: "/path/to/template/for-edit",
		},
		{
			name:        "edit mode should use vi, if EDITOR environment variable is missing",
			editorVar:   "",
			expectedCmd: "vi",
			expectedArg: "/path/to/template/for-edit",
		},
		{
			name:        "edit mode should return errors, if something goes wrong",
			editorVar:   "",
			shouldFail:  true,
			expectedCmd: "vi",
			expectedArg: "/path/to/template/for-edit",
			expectError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockRunner := &MockCommandRunner{
				ShouldFail: tt.shouldFail,
			}

			ioh := &IOHandler{
				LookupEnv: func(key string) (string, bool) {
					if key == "EDITOR" && tt.editorVar != "" {
						return tt.editorVar, true
					}
					return "", false
				},
				CommandRunner: mockRunner,
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

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, but got: %v", err)
			}

			if mockRunner.CapturedName != tt.expectedCmd {
				t.Fatalf("expected command %q, but got %q", tt.expectedCmd, mockRunner.CapturedName)
			}

			if len(mockRunner.CapturedArgs) != 1 {
				t.Fatalf("expected 1 argument, but got %d", len(mockRunner.CapturedArgs))
			}

			if mockRunner.CapturedArgs[0] != tt.expectedArg {
				t.Fatalf("expected argument %q, but got %q", tt.expectedArg, mockRunner.CapturedArgs[0])
			}
		})
	}
}

func TestListTemplates(t *testing.T) {
	var writer bytes.Buffer
	ioh := &IOHandler{
		Writer: &writer,
	}

	testCases := []struct {
		name          string
		storage       map[string]TemplateFile
		expect        string
		listTemplates bool
	}{
		{
			name: "happy path",
			storage: map[string]TemplateFile{
				"first": {
					Name:        "first",
					DefaultName: "",
					Path:        "/path/to/template/first",
				},
				"second": {
					Name:        "second",
					DefaultName: "parent",
					Path:        "/path/to/template/parent/second",
				},
			},
			expect: "first\nsecond - parent\n",
		},
		{
			name:    "empty storage",
			storage: map[string]TemplateFile{},
			expect:  "",
		},
		{
			name: "list templates - happy path",
			storage: map[string]TemplateFile{
				"first": {
					Name:        "first",
					DefaultName: "",
					Path:        "/path/to/template/first",
				},
				"second": {
					Name:        "second",
					DefaultName: "parent",
					Path:        "/path/to/template/parent/second",
				},
			},
			expect:        "first\nsecond\n",
			listTemplates: true,
		},
		{
			name:          "list templates - empty storage",
			storage:       map[string]TemplateFile{},
			expect:        "",
			listTemplates: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			defer writer.Reset()
			cliState := CliState{
				storage:       &Storage{templates: tt.storage},
				ioh:           ioh,
				listTemplates: tt.listTemplates,
			}

			err := cliState.Run()

			if err != nil {
				t.Fatalf("expected no error, but got: %v", err)
			}

			result := writer.String()
			if result != tt.expect {
				t.Fatalf("wrong output format, expected:\n%v\nbut got:\n%v\n", tt.expect, result)
			}
		})
	}
}

func TestTemplateOutput(t *testing.T) {
	var writer bytes.Buffer
	ioh := &IOHandler{
		Writer: &writer,
		LookupEnv: func(key string) (string, bool) {
			return key, true
		},
		ReadFile: func(name string) ([]byte, error) {
			return []byte(name + ": {VAR}\n"), nil
		},
		Create: func(name string) (OutputFile, error) {
			writer.Write([]byte(name))
			writer.Write([]byte("\n\n"))
			return StdoutInstance(&writer), nil
		},
	}

	testCases := []struct {
		name           string
		storage        map[string]TemplateFile
		outputFileName string
		defaultName    bool
		expect         string
	}{
		{
			name: "happy path",
			storage: map[string]TemplateFile{
				"first": {
					Name:        "first",
					DefaultName: "",
					Path:        "/path/to/template/first",
				},
			},
			expect: "/path/to/template/first: VAR\n",
		},
		{
			name: "template with default name but -d not set",
			storage: map[string]TemplateFile{
				"first": {
					Name:        "first",
					DefaultName: "default",
					Path:        "/path/to/template/first",
				},
			},
			expect: "/path/to/template/first: VAR\n",
		},
		{
			name: "output into file",
			storage: map[string]TemplateFile{
				"first": {
					Name:        "first",
					DefaultName: "",
					Path:        "/path/to/template/first",
				},
			},
			outputFileName: "second",
			expect:         "second\n\n/path/to/template/first: VAR\n",
		},
		{
			name: "output into file template with default name but -d not set",
			storage: map[string]TemplateFile{
				"first": {
					Name:        "first",
					DefaultName: "default",
					Path:        "/path/to/template/first",
				},
			},
			outputFileName: "second",
			expect:         "second\n\n/path/to/template/first: VAR\n",
		},
		{
			name: "output into file with default name",
			storage: map[string]TemplateFile{
				"first": {
					Name:        "first",
					DefaultName: "default",
					Path:        "/path/to/template/first",
				},
			},
			defaultName: true,
			expect:      "default\n\n/path/to/template/first: VAR\n",
		},
		{
			name: "a few templates",
			storage: map[string]TemplateFile{
				"first": {
					Name:        "first",
					DefaultName: "default",
					Path:        "/path/to/template/first",
				},
				"second": {
					Name:        "second",
					DefaultName: "",
					Path:        "/path/to/template/second",
				},
			},
			expect: "/path/to/template/first: VAR\n/path/to/template/second: VAR\n",
		},
		{
			name: "a few templates with output into file",
			storage: map[string]TemplateFile{
				"first": {
					Name:        "first",
					DefaultName: "default",
					Path:        "/path/to/template/first",
				},
				"second": {
					Name:        "second",
					DefaultName: "",
					Path:        "/path/to/template/second",
				},
			},
			outputFileName: "output",
			expect:         "output\n\n/path/to/template/first: VAR\noutput\n\n/path/to/template/second: VAR\n",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			defer writer.Reset()
			names := slices.Collect(maps.Keys(tt.storage))
			slices.Sort(names)
			cliState := CliState{
				templateNames:  names,
				storage:        &Storage{templates: tt.storage},
				ioh:            ioh,
				outputFileName: tt.outputFileName,
				defaultName:    tt.defaultName,
			}

			err := cliState.Run()

			if err != nil {
				t.Fatalf("expected no error, but got: %v", err)
			}

			result := writer.String()
			if result != tt.expect {
				t.Fatalf("wrong output format, expected:\n%v\nbut got:\n%v\n", tt.expect, result)
			}
		})
	}
}
