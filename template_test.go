package main

import (
	"slices"
	"testing"
)

func TestFindVariables(t *testing.T) {
	testCases := []struct {
		name      string
		content   string
		variables []string
	}{
		{
			name:      "empty string",
			content:   "",
			variables: []string{},
		},
		{
			name:      "without variables",
			content:   "just\nsome text",
			variables: []string{},
		},
		{
			name:      "single variable",
			content:   "variable {A}",
			variables: []string{"A"},
		},
		{
			name:      "same variable twice",
			content:   "variable {A} and {A}",
			variables: []string{"A"},
		},
		{
			name:      "empty variable",
			content:   "just\n\tsome text {} another text",
			variables: []string{},
		},
		{
			name:      "line break in variable name",
			content:   "line break {A\n}",
			variables: []string{},
		},
		{
			name:      "variable escaping",
			content:   "escaping \\{A} {B}",
			variables: []string{"B"},
		},
		{
			name:      "escaping without variable",
			content:   "fake escaping \\[A] {B}",
			variables: []string{"B"},
		},
		{
			name:      "few variables",
			content:   "one {A}, two{A B}, three {A B C}",
			variables: []string{"A", "A B", "A B C"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			template := NewTemplate("", []byte(tt.content))
			if !slices.Equal(template.Variables, tt.variables) {
				t.Fatalf("We should get %#v, but got %#v", tt.variables, template.Variables)
			}
		})
	}
}

func TestFillTemplate(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		result  string
		values  map[string]string
	}{
		{
			name:    "without variables",
			content: "Test\ntest",
			result:  "Test\ntest",
			values:  map[string]string{},
		},
		{
			name:    "empty variable",
			content: "Test {A}",
			result:  "Test ",
			values:  map[string]string{},
		},
		{
			name:    "empty variable with values",
			content: "Test {C}",
			result:  "Test ",
			values: map[string]string{
				"A": "A",
				"B": "B",
			},
		},
		{
			name:    "one variable",
			content: "Test {A}",
			result:  "Test Var",
			values: map[string]string{
				"A": "Var",
			},
		},
		{
			name:    "one variable and some text",
			content: "Test {A} another text",
			result:  "Test Var another text",
			values: map[string]string{
				"A": "Var",
			},
		},
		{
			name:    "two variables",
			content: "Test {A} rest {B}",
			result:  "Test 1 rest 2\n3",
			values: map[string]string{
				"A": "1",
				"B": "2\n3",
			},
		},
		{
			name:    "new line in a name",
			content: "Test {A\n} rest {B} and {B C}",
			result:  "Test {A\n} rest SOME TEXT and ANOTHER TEXT",
			values: map[string]string{
				"A":   "WITHOUT",
				"B":   "SOME TEXT",
				"B C": "ANOTHER TEXT",
			},
		},
		{
			name:    "some variables are empty",
			content: "Test {A} rest {B}!",
			result:  "Test  rest SOME TEXT!",
			values: map[string]string{
				"B": "SOME TEXT",
			},
		},
		{
			name:    "variable escaping",
			content: "Test \\{A} rest {B}!",
			result:  "Test {A} rest SOME TEXT!",
			values: map[string]string{
				"A": "WITHOUT",
				"B": "SOME TEXT",
			},
		},
		{
			name:    "escaping without variable",
			content: "Test \\[A] rest {B}!",
			result:  "Test \\[A] rest SOME TEXT!",
			values: map[string]string{
				"A": "WITHOUT",
				"B": "SOME TEXT",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			template := NewTemplate("", []byte(tt.content))
			result := template.fillTemplate(tt.values)
			if result != tt.result {
				t.Fatalf("We should get %#v, but got %#v", tt.result, result)
			}
		})
	}
}
