package parser

import (
	"slices"
	"testing"
)

type VariableTestCase struct {
	content   []byte
	variables []string
}

func NewVariableTestCase(content string, vars ...string) *VariableTestCase {
	testCase := new(VariableTestCase)
	testCase.content = []byte(content)
	testCase.variables = vars
	return testCase
}

func TestFindVariables(t *testing.T) {
	testCases := []*VariableTestCase{
		NewVariableTestCase(""),
		NewVariableTestCase("without variables"),
		NewVariableTestCase("variable {A}", "A"),
		NewVariableTestCase("variable {A} and {A}", "A"),
		NewVariableTestCase("empty variable {}"),
		NewVariableTestCase("line break {A\n}"),
		NewVariableTestCase("escaping \\{A}"),
		NewVariableTestCase("fake escaping \\[A]"),
		NewVariableTestCase("few variables {A}, {A B}, {A B C}", "A", "A B", "A B C"),
	}

	for _, testCase := range testCases {
		vars := FindVariables(testCase.content)
		if !slices.Equal(vars, testCase.variables) {
			t.Fatalf("We should get %#v, but got %#v", testCase.variables, vars)
		}
	}
}
