package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// StartSpinner запускает спиннер и возвращает функцию остановки.
// Гарантирует, что stop можно вызывать несколько раз без паники.
func StartSpinner() func() {
	stopCh := make(chan struct{})
	var once sync.Once

	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-stopCh:
				return
			default:
				fmt.Printf("\r %s%s%s thinking... ", ColorRed, frames[i], ColorReset)
				time.Sleep(80 * time.Millisecond)
				i = (i + 1) % len(frames)
			}
		}
	}()

	return func() {
		once.Do(func() {
			close(stopCh)
			// Очищаем строку спиннера
			fmt.Print("\r" + strings.Repeat(" ", 45) + "\r")
		})
	}
}
