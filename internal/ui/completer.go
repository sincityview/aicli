package ui

import (
	"strings"
	"sync"
)

type Completer struct{}

var (
	modelCache      []string
	modelCacheMutex sync.RWMutex
)

func (c Completer) Do(line []rune, pos int) ([][]rune, int) {
	str := string(line[:pos])

	// Статические команды
	static := []string{
		"/help ", "/status ", "/model ", "/provider ",
		"/chat ", "/clear ", "/read ", "/save ", "/exit ",
	}
	for _, cmd := range static {
		if strings.HasPrefix(cmd, str) {
			return [][]rune{[]rune(strings.TrimPrefix(cmd, str))}, len(line)
		}
	}

	return nil, len(line)
}
