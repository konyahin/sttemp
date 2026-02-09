package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestAskForValueStderrOutput(t *testing.T) {
	var writer bytes.Buffer
	reader := strings.NewReader("value")
	ioh := &IOHandler{
		Stdin:  reader,
		Stderr: &writer,
	}

	_, _ = ioh.askForValue("VAR")

	value := writer.String()
	expect := "Enter value for VAR: "
	if value != expect {
		t.Fatalf("expected: %v,\nbut got: %v.", expect, value)
	}
}

func TestAskForValue(t *testing.T) {
	testCases := []struct {
		name      string
		variable  string
		input     string
		expect    string
		expectErr error
	}{
		{
			name:     "happy path",
			variable: "VAR",
			input:    "VALUE\n",
			expect:   "VALUE",
		},
		{
			name:     "no trim for whitespaces",
			variable: "VAR",
			input:    " VALUE \n",
			expect:   " VALUE ",
		},
		{
			name:      "user cancel input",
			variable:  "VAR",
			input:     " VALUE ",
			expectErr: io.EOF,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var writer bytes.Buffer
			reader := strings.NewReader(tt.input)
			ioh := &IOHandler{
				Stdin:  reader,
				Stderr: &writer,
			}

			value, err := ioh.askForValue(tt.variable)
			if err != tt.expectErr {
				t.Fatalf("expected error %v,\nbut got: %v.", tt.expectErr, err)
			}

			if value != tt.expect {
				t.Fatalf("expected: %v,\nbut got: %v.", tt.expect, value)
			}
		})
	}
}
