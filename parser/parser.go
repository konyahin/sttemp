package parser

import (
	"iter"
	"maps"
	"slices"
)

type tokenType byte

const (
	Text tokenType = iota
	Variable
)

type parsingState byte

const (
	OutsideVar parsingState = iota
	InsideVar
	EscapingBracket
)

func tokens(content []byte) iter.Seq2[tokenType, []byte] {
	return func(yield func(tokenType, []byte) bool) {
		state := OutsideVar
		oldIdx := 0

		for i, c := range content {
			switch {
			case state == OutsideVar && c == '{':
				state = InsideVar
				if !yield(Text, content[oldIdx:i]) {
					return
				}
			case state == OutsideVar && c == '\\':
				state = EscapingBracket
				if !yield(Text, content[oldIdx:i]) {
					return
				}
			case state == EscapingBracket && c == '{':
				state = OutsideVar
			case state == EscapingBracket:
				state = OutsideVar
				continue
			case state == InsideVar && c == '}':
				state = OutsideVar
				if !yield(Variable, content[oldIdx+1:i]) {
					return
				}
			case state == InsideVar && c == '\n':
				state = OutsideVar
				if !yield(Text, content[oldIdx:i]) {
					return
				}
			default:
				continue
			}
			oldIdx = i
		}
	}
}

func FindVariables(content []byte) []string {

	vars := make(map[string]struct{})

	for token, content := range tokens(content) {
		if token == Variable && len(content) > 0 {
			vars[string(content)] = struct{}{}
		}
	}

	result := slices.Collect(maps.Keys(vars))
	slices.Sort(result)
	return result
}
