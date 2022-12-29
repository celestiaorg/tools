package tools

import (
	"fmt"
	"time"
)

func WaitCountdown(s time.Duration, msg string) {
	if msg != "" {
		msg += ", "
	}
	for ; s >= 0; s-- {
		time.Sleep(time.Second)
		fmt.Printf("\r%swaiting %3ds...", msg, s)
	}
}
