package cli

import (
	"bufio"
	"os"
	"strings"
)

func ConsolePrompt() string {
	r := bufio.NewReader(os.Stdin)
	input, _ := r.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input))
}
