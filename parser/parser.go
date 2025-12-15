package parser

import (
	"iter"
	"maps"
	"slices"
	"strings"
)

type Token struct {
	Type TokenType
	Content []byte
}

type TokenType byte

const (
	Text TokenType = iota
	Variable
)

type parsingState byte

const (
	OutsideVar parsingState = iota
	InsideVar
	EscapingBracket
)

func Tokens(content []byte) iter.Seq[Token] {
	return func(yield func(Token) bool) {
		state := OutsideVar
		oldIdx := 0

		for i, c := range content {
			switch {
			case state == OutsideVar && c == '{':
				state = InsideVar
				if !yield(Token{Text, content[oldIdx:i]}) {
					return
				}
			case state == OutsideVar && c == '\\':
				state = EscapingBracket
				if !yield(Token{Text, content[oldIdx:i]}) {
					return
				}
			case state == EscapingBracket && c == '{':
				state = OutsideVar
			case state == EscapingBracket:
				state = OutsideVar
				continue
			case state == InsideVar && c == '}':
				state = OutsideVar
				if !yield(Token{Variable, content[oldIdx+1:i]}) {
					return
				}
				oldIdx = i + 1
				continue
			case state == InsideVar && c == '\n':
				state = OutsideVar
				if !yield(Token{Text, content[oldIdx:i]}) {
					return
				}
			default:
				continue
			}
			oldIdx = i
		}

		if oldIdx < len(content) {
			yield(Token{Text, content[oldIdx:]})
		}
	}
}

func FindVariables(content []byte) []string {

	vars := make(map[string]struct{})

	for token := range Tokens(content) {
		if token.Type == Variable && len(token.Content) > 0 {
			vars[string(token.Content)] = struct{}{}
		}
	}

	result := slices.Collect(maps.Keys(vars))
	slices.Sort(result)
	return result
}

func FillTemplate(content []byte, values map[string]string) string {
	var sb strings.Builder
	for token := range Tokens(content) {
		switch {
		case token.Type == Text:
			sb.Write(token.Content)
		case token.Type == Variable && len(token.Content) > 0:
			val, ok := values[string(token.Content)]
			if ok {
				sb.WriteString(val)
			}
		}
	}
	return sb.String()
}
