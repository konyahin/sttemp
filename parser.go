package main

type parsingState int8

const (
	outsideVar parsingState = iota
	insideVar
	escapingBracket
)

func findVariables(content []byte) []string {
	state := outsideVar

	var vars []string
	var buf []byte

	for _, c := range content {
		switch {
		case state == outsideVar && c == '{':
			state = insideVar
		case state == outsideVar && c == '\\':
			state = escapingBracket
		case state == escapingBracket:
			state = outsideVar
		case state == insideVar && c != '}' && c != '\n':
			buf = append(buf, c)
		case state == insideVar && c == '}':
			state = outsideVar
			vars = append(vars, string(buf))
			buf = nil
		case state == insideVar && c == '\n':
			state = outsideVar
			buf = nil
		}
	}

	return vars
}
