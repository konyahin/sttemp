package parser

import (
	"slices"
	"testing"
)

type VariableTestCase struct {
	content   []byte
	variables []string
}

func newVariableTestCase(content string, vars ...string) *VariableTestCase {
	testCase := new(VariableTestCase)
	testCase.content = []byte(content)
	testCase.variables = vars
	return testCase
}

type FillTemplateTestCase struct {
	content string
	result  string
	values  map[string]string
}

func newFillTemplateTestCase(content string, result string, values  map[string]string) *FillTemplateTestCase {
	testCase := new(FillTemplateTestCase)
	testCase.content = content
	testCase.result = result
	testCase.values = values
	return testCase
}

func TestFindVariables(t *testing.T) {
	testCases := []*VariableTestCase{
		newVariableTestCase(""),
		newVariableTestCase("without variables"),
		newVariableTestCase("variable {A}", "A"),
		newVariableTestCase("variable {A} and {A}", "A"),
		newVariableTestCase("empty variable {}"),
		newVariableTestCase("line break {A\n}"),
		newVariableTestCase("escaping \\{A}"),
		newVariableTestCase("fake escaping \\[A]"),
		newVariableTestCase("few variables {A}, {A B}, {A B C}", "A", "A B", "A B C"),
	}

	for _, testCase := range testCases {
		vars := FindVariables(testCase.content)
		if !slices.Equal(vars, testCase.variables) {
			t.Fatalf("We should get %#v, but got %#v", testCase.variables, vars)
		}
	}
}

func TestFillTemplate(t *testing.T) {
	testCases := []*FillTemplateTestCase{
		newFillTemplateTestCase("Test", "Test", nil),
		newFillTemplateTestCase("Test {A}", "Test ", nil),
		newFillTemplateTestCase("Test", "Test", map[string]string{
			"A": "B",
		}),
		newFillTemplateTestCase("Test {A}", "Test B", map[string]string{
			"A": "B",
		}),
		newFillTemplateTestCase("Test {A} rest", "Test B rest", map[string]string{
			"A": "B",
			"B": "C",
		}),
		newFillTemplateTestCase("Test {A} rest {B}", "Test 1 rest 2", map[string]string{
			"A": "1",
			"B": "2",
		}),
		newFillTemplateTestCase("Test {A\n} rest {B} and {B C}", "Test {A\n} rest SOME TEXT and ANOTHER TEXT", map[string]string{
			"A": "WITHOUT",
			"B": "SOME TEXT",
			"B C": "ANOTHER TEXT",
		}),
		newFillTemplateTestCase("Test {A} rest {B}!", "Test  rest SOME TEXT!", map[string]string{
			"B": "SOME TEXT",
		}),
	}

	for _, testCase := range testCases {
		result := FillTemplate([]byte(testCase.content), testCase.values)
		if result != testCase.result {
			t.Errorf("Test case: %#v", testCase)
			t.Fatalf("We should get %#v, but got %#v", testCase.result, result)
		}
	}
}
