package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ConsolePrompt() string {
	r := bufio.NewReader(os.Stdin)
	input, _ := r.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input))
}

func Ask(prompt, no string) bool {
	fmt.Println(prompt)
	answer := ConsolePrompt()
	switch answer {
	case no:
		return false // caller gracefully exits
	}
	return true
}
