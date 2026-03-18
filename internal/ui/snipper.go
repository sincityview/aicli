package ui

import (
	"fmt"
	"strings"
	"time"
)

func StartSpinner() func() {
	stopCh := make(chan struct{})
	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-stopCh:
				fmt.Print("\r" + strings.Repeat(" ", 45) + "\r")
				return
			default:
				fmt.Printf("\r %s%s%s thinking... ", ColorRed, frames[i], ColorReset)
				time.Sleep(80 * time.Millisecond)
				i = (i + 1) % len(frames)
			}
		}
	}()
	return func() { close(stopCh) }
}
