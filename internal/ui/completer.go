package ui

import "strings"

type Completer struct{}

func (c Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	str := string(line[:pos])
	commands := []string{
		"/help", "/list", "/new", "/clear", "/models", "/exit",
		"/model ", "/open ", "/delete ", "/rename ", "/read ", "/save ",
	}

	for _, cmd := range commands {
		if strings.HasPrefix(cmd, str) {
			newLine = append(newLine, []rune(strings.TrimPrefix(cmd, str)))
		}
	}
	return newLine, len(line)
}
