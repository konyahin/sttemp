package main

type Token struct {
	Type    TokenType
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

func tokens(content []byte) []Token {
	tokens := make([]Token, 0)

	state := OutsideVar
	oldIdx := 0

	for i, c := range content {
		switch {
		case state == OutsideVar && c == '{':
			state = InsideVar
			tokens = append(tokens, Token{Text, content[oldIdx:i]})
		case state == OutsideVar && c == '\\':
			state = EscapingBracket
			tokens = append(tokens, Token{Text, content[oldIdx:i]})
		case state == EscapingBracket && c == '{':
			state = OutsideVar
		case state == EscapingBracket:
			state = OutsideVar
			continue
		case state == InsideVar && c == '}':
			state = OutsideVar
			tokens = append(tokens, Token{Variable, content[oldIdx+1 : i]})
			oldIdx = i + 1
			continue
		case state == InsideVar && c == '\n':
			state = OutsideVar
			tokens = append(tokens, Token{Text, content[oldIdx:i]})
		default:
			continue
		}
		oldIdx = i
	}

	if oldIdx < len(content) {
		tokens = append(tokens, Token{Text, content[oldIdx:]})
	}

	return tokens
}
