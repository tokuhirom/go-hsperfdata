package attach

import (
	"regexp"
	"strconv"
	"strings"
)

type JavaThread struct {
	name   string
	state  string
	stacks []*StackItem
}

type StackItem struct {
	method string
	file   string
	line   int
}

func (self *JavaThread) GetName() string {
	return self.name
}
func (self *JavaThread) GetState() string {
	return self.state
}
func (self *JavaThread) GetStack() []*StackItem {
	return self.stacks
}
func (self *StackItem) GetMethod() string {
	return self.method
}
func (self *StackItem) GetFile() string {
	return self.file
}
func (self *StackItem) GetLine() int {
	return self.line
}

var header_re = regexp.MustCompile(`^"([^"]+)"`)
var thread_state_re = regexp.MustCompile(`java.lang.Thread.State: ([A-Z_]+)(?: \((.*)\))?`)
var stack_at_re = regexp.MustCompile(`\s*at ([^(]+)\((.*?)(?::(\d+))?\)`)

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
			file := matches[2]
			var line int
			if len(matches) == 4 {
				l, err := strconv.Atoi(matches[3])
				if err != nil {
					line = -1
				} else {
					line = l
				}
			} else {
				line = -1
			}
			stacks = append(stacks, &StackItem{matches[1], file, line})
		} else if len(name) > 0 {
			threads = append(threads, &JavaThread{name, state, stacks})
			name = ""
			state = ""
			stacks = stacks[:0]
		}
	}
	return threads, nil
}
