package main

import (
	"slices"
	"testing"
)

type VariableTestCase struct {
	content   []byte
	variables []string
}

func NewVariableTestCase(vars ...string) *VariableTestCase {
	testCase := new(VariableTestCase)
	testCase.content = []byte(vars[0])
	testCase.variables = vars[1:]
	return testCase
}

func TestFindVariables(t *testing.T) {
	testCases := []*VariableTestCase{
		NewVariableTestCase(""),
		NewVariableTestCase("without variables"),
		NewVariableTestCase("variable {A}", "A"),
		NewVariableTestCase("empty variable {}"),
		NewVariableTestCase("line break {A\n}"),
		NewVariableTestCase("escaping \\{A}"),
		NewVariableTestCase("fake escaping \\[A]"),
		NewVariableTestCase("few variables {A}, {A B}, {A B C}", "A", "A B", "A B C"),
	}

	for _, testCase := range testCases {
		vars := findVariables(testCase.content)
		if !slices.Equal(vars, testCase.variables) {
			t.Fatalf("We should get %#v, but got %#v", testCase.variables, vars)
		}
	}
}
