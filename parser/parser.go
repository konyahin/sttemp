package parser

import (
	"iter"
	"maps"
	"slices"
	"strings"
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
				oldIdx = i + 1
				continue
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

		if oldIdx < len(content) {
			yield(Text, content[oldIdx:])
		}
	}
}

func FindVariables(content []byte) []string {

	vars := make(map[string]struct{})

	for token, value := range tokens(content) {
		if token == Variable && len(value) > 0 {
			vars[string(value)] = struct{}{}
		}
	}

	result := slices.Collect(maps.Keys(vars))
	slices.Sort(result)
	return result
}

func FillTemplate(content []byte, values map[string]string) string {
	var sb strings.Builder
	for token, value := range tokens(content) {
		switch {
		case token == Text:
			sb.Write(value)
		case token == Variable && len(value) > 0:
			val, ok := values[string(value)]
			if ok {
				sb.WriteString(val)
			}
		}
	}
	return sb.String()
}
