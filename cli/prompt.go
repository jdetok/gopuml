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

func AskBool(prompt, no string) bool {
	fmt.Print(prompt)
	answer := ConsolePrompt()
	switch answer {
	case no:
		return false // caller gracefully exits
	}
	return true
}

func AskStr(prompt string) string {
	fmt.Print(prompt)
	return ConsolePrompt()
}
