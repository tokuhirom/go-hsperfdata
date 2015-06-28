package attach

import (
	"regexp"
	"strings"
)

type JavaThread struct {
	name   string
	state  string
	stacks []*StackItem
}

type StackItem struct {
	method string
}

var header_re = regexp.MustCompile(`^"([^"]+)"`)
var thread_state_re = regexp.MustCompile(`java.lang.Thread.State: ([A-Z_]+)(?: \((.*)\))?`)
var stack_at_re = regexp.MustCompile(`\s*at ([^(]+)\(`)

func ParseStack(stack string) ([]*JavaThread, error) {

	threads := make([]*JavaThread, 0)
	lines := strings.Split(stack, "\n")
	name := ""
	state := ""
	stacks := make([]*StackItem, 0)
	for _, line := range lines {
		if matches := header_re.FindStringSubmatch(line); len(matches) == 2 {
			name = matches[1]
		} else if matches = thread_state_re.FindStringSubmatch(line); len(matches) >= 2 {
			state = matches[1]
		} else if matches = stack_at_re.FindStringSubmatch(line); len(matches) >= 1 {
			stacks = append(stacks, &StackItem{matches[1]})
		} else if len(name) > 0 {
			threads = append(threads, &JavaThread{name, state, stacks})
			name = ""
			state = ""
			stacks = stacks[:0]
		}
	}
	return threads, nil
}
